This repo is intended to gather feedback and opinions on few good practices and techniques that could be ***helpful to craft better web servers in Go***.

Parts of the minimal web framework have been extracted and adapted from the excellent [ardanlabs/service](https://github.com/ardanlabs/service) repo.

The technique I'm trying to fit into the web server is the `opaque errors` one, by Dave Cheney.
(And yes, I'm a big fan of both Dave and Bill :smiley:)

The two main concepts on error handling that led me to this implementation are:
- handle errors once, in a centralized fashion using dedicated middleware.
- leverage opaque errors to decouple errors definitions from the web framework.

**Refs**:
- https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
- https://www.ardanlabs.com/blog/2017/05/design-philosophy-on-logging.html
- https://github.com/ardanlabs/service

