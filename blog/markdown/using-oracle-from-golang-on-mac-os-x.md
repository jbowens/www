# Using Oracle from Go on Mac OS X

I recently found myself needing to use Oracle from an application written in Go. I have zero experience with Oracle, and I found it pretty difficult to get a working development environment setup on my Macbook. Here's a walkthough of everything I did to get a workable environment:

 * Download and install [VirtualBox](https://www.virtualbox.org/wiki/VirtualBox) if you don't already have it.
 * Download the [Oracle DB Developer VM](http://www.oracle.com/technetwork/database/enterprise-edition/databaseappdev-vm-161299.html) (warning: it's almost 8gb).
 * Download the "Basic", "SDK" and "SQL\*Plus" packages of the [Oracle Instant Client](http://www.oracle.com/technetwork/topics/intel-macsoft-096467.html).

## Set up the Oracle DB Developer VM

First, in Finder, double click on the `.ova` file you downloaded. Follow VirtualBox's prompts to import the virtual machine. It may take a couple minutes. Once it's imported, start up the VM. In the VM's open Terminal shell, run `sqlplus`. When prompted for user-name and password, provide `system` and `oracle` respectively. Then create a privileged Oracle user for testing. For example:

```
CREATE USER jackson IDENTIFIED BY password;
GRANT CONNECT, RESOURCE, DBA, CREATE SESSION, UNLIMITED TABLESPACE TO jackson;
```

Back in your Mac's terminal, verify that you can connect to your VM. Try SSHing into your VM as the `oracle` user with password `oracle`. The VM's sshd should be listening on localhost port 2222.

```
$ ssh -p 2222 oracle@localhost
oracle@localhost's password:
Last login: Wed Apr  5 18:07:15 2017
[oracle@vbgeneric ~]$
```

If that doesn't work, you may need to play with VMWare's port forwarding settings. If you run into issues setting up the VM, [these instructions](http://www.thatjeffsmith.com/archive/2014/02/introducing-the-otn-developer-day-database-12c-virtualbox-image/) may be helpful.

## Set up the Oracle Instant Client

The Oracle Instant Client provides the Oracle Client Interface (OCI) dynamic library required by all Go drivers (and drivers for most languages). Make a directory somewhere for your Instant Client to live. I used `~/oracle`. Unzip all three zips into the directory.

```
unzip -d ~/oracle/ 'instantclient-*.zip'
```

The directory should look something like:
```
oracle/
└── instantclient_12_1
    ├── BASIC_README
    ├── SQLPLUS_README
    .
    .
    .
    ├── libclntsh.dylib.12.1
    ├── libclntshcore.dylib.12.1
    .
    .
    .
    ├── sdk
    │   ├── SDK_README
        │   ├── include
    │   │   ├── ldap.h
    │   │   ├── nzerror.h
    │   │   ├── nzt.h
    │   │   ├── occi.h
    sqlplus
    .
    .
    .
```

Next, you'll want to update a few environment variables to point to this new directory. You can do this just in your current session or in your `.bash_profile`:

```
export PATH=$PATH:~/oracle/instantclient_12_1/
export DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH:~/oracle/instantclient_12_1/
```

Use the `sqlplus` command to test connecting to Oracle inside your VM. You'll need to provide the DSN for the database as an argument to sqlplus. The only values you should need to change are the username and password.

```
sqlplus jackson/password@127.0.0.1:1521/orcl
```

Once you're connected, you can issue any DDL statements you want to setup your test environment. If you want you can also seed data.

```
SQL> CREATE TABLE juices (
  2    name varchar2(100)
  3  );

Table created.

SQL> INSERT INTO juices (name) VALUES('orange juice');

1 row created.

SQL> COMMIT;

Commit complete.
```

## Set up the Go driver

Before we can install any of the Go drivers for Oracle, we need to ensure the compiler will be able to find the Oracle Instant Client. Create a new file `/usr/local/lib/pkgconfig/oci8.pc` containing the following text, replacing the `prefix` value with the location of your Instant Client:

```
prefix=/Users/jackson/oracle/instantclient_12_1
version=12.1
libdir=${prefix}
includedir=${prefix}/sdk/include

Name: oci8
Description: Oracle database engine
Version: ${version}
Libs: -L${libdir} -lclntsh
Libs.private:
Cflags: -I${includedir}
```

Try installing a Go driver such as [go-oci8](https://github.com/mattn/go-oci8) or [ora](https://github.com/rana/ora):
```
go get github.com/mattn/go-oci8
```

If you encounter a "`ld: library not found for -lclntsh`" error, the version suffixes on some the libraries' filenames might be confusing the linker.

<details>
```
$ ls -l *.12.1
-rwxrwxrwx@ 1 jackson  staff  67437336 Jun  8  2016 libclntsh.dylib.12.1
-rwxrwxrwx@ 1 jackson  staff   4532196 Jun  8  2016 libclntshcore.dylib.12.1
-rwxrwxrwx@ 1 jackson  staff   1483956 Jun  8  2016 libocci.dylib.12.1
-rwxrwxrwx@ 1 jackson  staff   5415256 Jun  8  2016 libsqora.dylib.12.1
```

Try copying each one of these files to an identical file without the `.12.1` suffix and try again.
```
cp libclntsh.dylib.12.1 libclntsh.dylib
cp libclntshcore.dylib.12.1 libclntshcore.dylib
cp libocci.dylib.12.1 libocci.dylib
cp libsqora.dylib.12.1 libsqora.dylib
```
</details>

## Try it!

Try compiling and running a simple Go program using the Oracle driver. An example program querying the `juices` table I created up above:

```
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-oci8"
)

const dsn = `jackson/password@127.0.0.1:1521/orcl`

func main() {
	db, err := sql.Open("oci8", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM juices")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		fmt.Println(name)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
```

If when running your program you get an `signal: killed` error like I did, you might be hitting [golang/go#19734](https://github.com/golang/go/issues/19734). Thankfully, there's a workaround using `go build -ldflags -s`.

Enjoy your new life as an Oracle developer. ;)
