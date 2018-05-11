[![Docker Repository on Quay](https://quay.io/repository/project-prismatica/prismatica-report-render/status "Docker Repository on Quay")](https://quay.io/repository/project-prismatica/prismatica-report-render)
# Report Renderer

Report Renderer is a gRPC service which will render reports given templates of
the reports and a configuration for remote data sources. It builds on
[pongo2](https://github.com/flosch/pongo2) with a few new functions for
accessing data used by the Project Prismatica ecosystem. Markdown support is
provided by [blackfriday](https://github.com/russross/blackfriday) through the
[pongo2-addons](https://github.com/flosch/pongo2-addons) package.

# Example Usage


## Basic Example
Template stored in ```examples/basic.tpl```
```
This is a basic template. The answer is {{ 41 | add:1 }}
```

Rendering:
```bash
$ prismatica_report_renderer render -s examples/basic.tpl
This is a basic template. The answer is 42
```

## Variable Example
Template stored in ```examples/var.tpl```
```
The variable provided was {{ asdf }}
```

Rendering:
```bash
$ prismatica_report_renderer render -s examples/var.tpl asdf=Something
The variable provided was Something
```

## As a service

While local rendering is useful for debugging, the power of the project is to
run this as a gRPC-exposed service as described in
```/prismatica_report_renderer.proto``` to allow clients to leverage data stored
in back-end mongo databases.

[Prismatica Infrastructure](https://github.com/Project-Prismatica/prismatica-infrastructure)
is the recommended method for running this as a service.

# XML Features

More documentation (here)[./blob/master/doc/xml.md]

# MongoDB Queries

More documentation (here)[./blob/master/doc/mongo.md]
