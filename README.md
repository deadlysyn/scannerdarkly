# Contents

- [Introduction](#route53-zone-auditor)
- [Usage](#usage)
- [Dependencies](#dependencies)
- [References](#references)

## Route53 Zone Auditor

Crawls records in hosted zones and checks for listening ports to help
prevent DNS hijacking. Stale DNS aliases may allow crackers to stand
something up at target addresses and masquerade as you.

Only aliases in public zones are audited by default since these pose
the most risk in cloud environments. You can optionally enable A and
AAAA. MX, NS, SOA, TXT and ACM-related CNAMEs are ignored.

While you can run local scans easily with [aws-vault](https://github.com/99designs/aws-vault),
the ideal place is a pipeline or EC2 instance within AWS. This is
because AWS generally blocks anything that looks like "scanning" --
the exception is when scanning your own resources. Be aware of
[the guidelines](https://aws.amazon.com/security/penetration-testing).
For Route53 zone enumeration in particular, the approach is to leverage
APIs vs crawling public DNS infrastructure. The latter violates AWS policy.

## Usage

Thanks to [viper](https://github.com/spf13/viper) and
[cobra](https://github.com/spf13/cobra), configuration can be
specified as environment variables, YAML configuration or
command line arguments.

By default all public hosted zones in the target account are
enumerated, but you can limit to specific zones.

```console
❯ aws-vault exec dev -- ./d -h
Scan for dark (stale) Route53 records

Usage:
  d [flags]

Flags:
  -a, --all             scan A/AAAA records (in addition to aliases)
  -c, --config string   config file (default "config.yml")
  -h, --help            help for d
  -j, --json            json output
  -n, --name string     report file name (default "report")
  -p, --ports strings   TCP ports to scan
  -t, --timeout int     port scan timeout (default 10)
  -z, --zones strings   zone ids to scan

❯ aws-vault exec dev -- ./d -r custom_report.csv
[status output trimmed]
Writing report: report

❯ cat custom_report.csv
Zone ID,Name,Type,Values (no open ports)
Zabc...,foo.bar.com,CNAME,somehost.somedomain
Zdef...,baz.bar.com,CNAME,anotherhost.somedomain

❯ aws-vault exec dev -- ./d -j -p "80 443" -z "Zabc... Zdef..." -t 5
{
        "Zabc...": [
                {
                        "Name": "foo.com",
                        "Type": "A",
                        "Alias": true,
                        "Values": [
                                "bar.az-1.elasticbeanstalk.com"
                        ],
                        "Active": [
                                "baz.az-1.elasticbeanstalk.com:80"
                        ]
                },
...
```

The environment prefix is `PKD_`.

```console
$ PKD_PORTS="443 4443 8443" PKD_ZONES="ZXXX..." aws-vault exec dev -- ./d
Processing zone ZXXX...
Skipping foo.domain.dev (NS)
Skipping foo.domain.dev (SOA)
Skipping _XXX.foo.domain.dev (ACM)
Scanning bar.az-1.elb.amazonaws.com:80... open.

$ cat report.csv
Zone ID,Name,Type,Results (no open ports)
ZXXX...,foo.domain.dev,Alias,bar.az-1.elb.amazonaws.com
...
```

## TODO

Parallel scanning would speed things up, but look more like a DDoS.
Alternate approach may be removing scanning entirely and feeding
lists of host:ports to an external tool like nmap.

- Parallel scanning
- More scan types (HEAD, version check, etc.)
- Scan more record types (External NS? Bogus MX?)

## Dependencies

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go-v2)

## References

- [SDK Getting Started](https://aws.github.io/aws-sdk-go-v2/docs/getting-started)
- [AWS Penetration Test Guidelines](https://aws.amazon.com/security/penetration-testing)
