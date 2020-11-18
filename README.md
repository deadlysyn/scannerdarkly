# Contents

- [Introduction](#route53-zone-auditor)
- [Usage](#usage)
  - [Example](#example)
- [Dependencies](#dependencies)
- [References](#references)

## Route53 Zone Auditor

Crawls records in hosted zones and checks for listening ports.

Since the idea is preventing DNS hijacking, only public zones are audited.
Any A, AAAA, CNAME or ALIAS records are port scanned. NS, SOA and ACM-related
records are ignored.

This is MVP, so only generates a CSV report. JSON output TBD.

## Usage

All configuration can be specified via environment:

- `PORTS`: Space delimited list of TCP ports to check (default: 80 443 8080 8443)
- `TIMEOUT`: Scanning timeout (includes DNS resolution, default: 10 seconds)
- `ZONES`: Space delimited list of hosted zone IDs to audit (default: all public zones in account)

### Example

```console
$ PORTS=80 ZONES=ZXXX... aws-vault exec dev -- go run . > report.csv
Processing zone ZXXX...
Skipping foo.domain.dev (NS)
Skipping foo.domain.dev (SOA)
Skipping _XXX.foo.domain.dev (ACM)
Scanning bar.region.elb.amazonaws.com:80... open.

$ cat report.csv
Zone ID,Name,Type,Results
ZXXX...,foo.domain.dev,Alias,bar.region.elb.amazonaws.com:80
...
```

## TODO

- Parallel scanning
- More scan types (HEAD, version check, etc.)
- Better arg/environment parsing
- JSON output (feed to other tools?)

## Dependencies

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go)

## References

- [AWS SDK API Doc](https://docs.aws.amazon.com/sdk-for-go/api)
- [AWS Penetration Test Guidelines](https://aws.amazon.com/security/penetration-testing)

