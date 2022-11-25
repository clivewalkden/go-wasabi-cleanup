# Cleanup old files on WasabiSys

This executable automatically cleans up old files outside the given compliance timeframes.

## Run
`go run cmd/main.go`

## Config file
This application requires a configuration file to be present in the user home directory or the directory the executable is being run from.

An example config file is included in this repository [wasabi-cleanup.sample.yml](./wasabi-cleanup.sample.yml)

This file contains the following:
```yaml
buckets:
  bucket-name: 90
  bucket-name-2: 180
  bucket-name-3: 365
connection:
  url: 'https://s3.us-central-1.wasabisys.com'
  region: us-central-1
  profile: wasabi
```
1. `buckets` are the names of the buckets you want to delete files from and the number of days back from today you want kept in the bucket.
2. `connection` the connection information pointing to the server with the files housed
   1. `url` the server access url
   2. `region` the region the server is located
   3. `profile` the AWS config and credentials profile used to connect

## Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/clivewalkden/go-wasabi-cleanup/tags) or [CHANGELOG.md](./CHANGELOG.md).


## Authors
* **Clive Walkden** - *Initial work* - [SOZO Design Ltd](https://github.com/sozo-design)

See also the list of [contributors](https://github.com/clivewalkden/go-wasabi-cleanup/contributors) who participated in this project.


## License
This project is licensed under the Proprietary License - see the [LICENSE](./LICENSE) file for details