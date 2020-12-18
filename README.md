# Contents

- [Introduction](#route53-zone-auditor)
- [Usage](#usage)
  - [Example](#example)
- [Dependencies](#dependencies)
- [References](#references)

## Route53 Zone Auditor

_NOTE: This is MVP, so only generates a CSV report. JSON output TBD. Other ideas?_

Crawls records in hosted zones and checks for listening ports. Only public
zones are audited. Any A, AAAA, CNAME or ALIAS records are port scanned.
MX, NS, SOA, TXT and ACM-related CNAMEs are ignored.

The idea is helping prevent DNS hijacking. If you have stale DNS records
in your zones, a would-be "hacker" (e.g. a bored cracker) can potentially
stand something up at the former address and masquerade as your domain.

While you can run local scans easily with aws-vault, the ideal place is
a pipeline or EC2 instance within AWS. This is because AWS generally
blocks anything that looks like "scanning" -- the exception is when
scanning your own resources. Be aware of
[the guidelines](https://aws.amazon.com/security/penetration-testing).
For Route53 zone enumeration in particular, the approach is to leverage
APIs vs crawling public DNS infrastructure. DO NOT attempt the latter
or you are in violation of AWS policy.

## Usage

All configuration can be specified via environment:

- `PORTS`: Space delimited list of TCP ports to check (default: 80 443)
- `TIMEOUT`: Scanning timeout (includes DNS resolution, default: 5 seconds)
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
- Scan more record types (External NS? Bogus MX?)
- Better arg/environment parsing
- JSON output (feed to other tools?)

## Dependencies

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go)

## References

- [AWS SDK API Doc](https://docs.aws.amazon.com/sdk-for-go/api)
- [AWS Penetration Test Guidelines](https://aws.amazon.com/security/penetration-testing)

