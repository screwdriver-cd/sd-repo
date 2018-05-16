# sd-repo
[![Build Status][status-image]][status-url] [![Open Issues][issues-image]][issues-url]

> Tool for executing the [Repo](https://source.android.com/setup/develop/repo) workflow for `getCheckoutCommand` in [screwdriver-scm-github](https://github.com/screwdriver-cd/scm-github/blob/master/index.js#L272)

## Usage

```bash
sd-repo -manifestUrl=git@github.com:org/manifestRepo.git/default.xml -sourceRepo=org/appRepo
```

## Testing

```bash
go test ./...
```

## License

Code licensed under the BSD 3-Clause license. See LICENSE file for terms.

[issues-image]: https://img.shields.io/github/issues/screwdriver-cd/screwdriver.svg
[issues-url]: https://github.com/screwdriver-cd/screwdriver/issues
[status-image]: https://cd.screwdriver.cd/pipelines/796/badge
[status-url]: https://cd.screwdriver.cd/pipelines/796
