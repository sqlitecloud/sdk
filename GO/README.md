# SQLite Cloud GO Client

## Get started

### Fetching the code:
```console
cd
mkdir test
cd test
git clone https://github.com/sqlitecloud/sdk
User: <github user name>
Password: <access-token or password>

```

### Setting up the development environment
```console
go env -w GO111MODULE=off
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH

```

### Install and compile the pre-requirements
First, you will have to install all the pre-requirements on your machine: `make install-prerequirements`. This will also patch the source code of the linenoise package...

### Building the test programs
If you want to build the Test programs: `make test`

### Building the CLI App
To build the CLI App (Warning: not fully functional, this is officially Step 1), you have to enter: `make cli`

### Build all at the same time:
If you want to do all at the same time: `make all`

## Documentation
If you want to see the Documentation: `make doc` - Warning: A browser window will open and display the documentation to you. The Documentation is updated live while coding. To stop the live mode, press CRTL-C on the command line.

## Development helpers
- Check files with gosec: `make checksec`
- Open the repo in github: `make github`.
- See changes: `make diff`
- Clean dependencies and precompiled code: `make clean`