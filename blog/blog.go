package blog

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/boltdb/bolt"
	"github.com/russross/blackfriday"
)

var (
	mu    sync.Mutex
	db    *bolt.DB
	posts map[string]Post = map[string]Post{}
)

var (
	bucketPosts = []byte("posts")
)

type byCreatedAt []Post

func (p byCreatedAt) Len() int           { return len(p) }
func (p byCreatedAt) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byCreatedAt) Less(i, j int) bool { return p[i].CreatedAt.After(p[j].CreatedAt) }

// Post represents an individual blog entry.
type Post struct {
	ID        string
	Markdown  []byte
	Plaintext string
	HTML      template.HTML
	Metadata
}

func (p Post) Title() string {
	if p.OverrideTitle != "" {
		return p.OverrideTitle
	}

	lines := strings.Split(p.Plaintext, "\n")
	title := strings.TrimSpace(lines[0])
	if len(title) > 0 {
		return title
	}

	t := strings.Join(strings.Split(p.ID, "-"), " ")
	runes := []rune(t)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (p Post) Snippet() string {
	paragraphs := strings.Split(p.Plaintext, "\n\n")
	if len(paragraphs) == 0 {
		return ""
	}
	lines := strings.SplitN(paragraphs[0], "\n", 2)
	if len(lines) > 1 {
		return strings.TrimSpace(lines[1])
	}
	if len(paragraphs[0]) > 0 {
		return strings.TrimSpace(paragraphs[0])
	}
	return p.Title()
}

// StructuredData returns JSON-LD json object describing the post.
func (p Post) StructuredData() template.HTML {
	// TODO(jackson): Cache the structured data output on the post?
	b, err := json.Marshal(map[string]interface{}{
		"@context": "http://schema.org",
		"@type":    "BlogPosting",
		"mainEntityOfPage": map[string]interface{}{
			"@type": "WebPage",
			"@id":   fmt.Sprintf("https://www.jbowens.org/p/%s", p.ID),
		},
		"headline":      p.Title(),
		"datePublished": p.CreatedAt,
		"dateModified":  p.UpdatedAt,
		"author": map[string]interface{}{
			"@type": "Person",
			"name":  "Jackson Owens",
		},
		"publisher": map[string]interface{}{
			"@type": "Person",
			"name":  "Jackson Owens",
		},
	})
	if err != nil {
		panic(err)
	}
	return template.HTML(b)
}

// Metadata represents metadata about a post. It's stored in boltdb, instead
// of on the file system with the markdown.
type Metadata struct {
	Hash          [32]byte
	OverrideTitle string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Posts returns all of the posts in the blog.
func Posts() []Post {
	mu.Lock()
	defer mu.Unlock()

	copiedPosts := make([]Post, 0, len(posts))
	for _, p := range posts {
		copiedPosts = append(copiedPosts, p)
	}

	sort.Sort(byCreatedAt(copiedPosts))
	return copiedPosts
}

// PostByID returns the post with the provided id, if it exists.
func PostByID(id string) (Post, bool) {
	mu.Lock()
	defer mu.Unlock()

	p, ok := posts[id]
	return p, ok
}

// Load takes a directory and loads all of the markdown files in the
// directory as posts.
func Load(dir string) error {
	var err error
	db, err = bolt.Open("data/blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Create the posts bucket if it doesn't already exist.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketPosts)
		return err
	})
	if err != nil {
		return err
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := filepath.Base(path)
		name = strings.TrimSuffix(name, ".md")

		p := Post{
			ID:        name,
			Markdown:  b,
			Plaintext: string(blackfriday.Markdown(b, snippetRenderer{}, 0)),
			HTML:      template.HTML(blackfriday.MarkdownCommon(b)),
			Metadata:  Metadata{},
		}
		markdownHash := sha256.New().Sum(b)
		copy(p.Metadata.Hash[:], markdownHash)

		err = lookupMetadata(&p)
		if err != nil {
			return err
		}

		posts[name] = p
		return nil
	})
	return err
}

func lookupMetadata(p *Post) error {
	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketPosts)
		v := b.Get([]byte(p.ID))

		var (
			m   Metadata
			err error
		)
		if v != nil {
			decoder := gob.NewDecoder(bytes.NewReader(v))
			err = decoder.Decode(&m)
			if err != nil {
				return err
			}
		}

		// We need to update the metadata and do a Put.
		if !bytes.Equal(p.Hash[:], m.Hash[:]) {
			m.Hash = p.Hash
			if m.CreatedAt.IsZero() {
				m.CreatedAt = time.Now()
			}
			m.UpdatedAt = time.Now()

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			err = encoder.Encode(m)
			if err != nil {
				return err
			}

			err = b.Put([]byte(p.ID), buf.Bytes())
			if err != nil {
				return err
			}
		}
		p.Metadata = m
		return err
	})
	return err
}
