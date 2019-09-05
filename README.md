# rd

This is `run-directory` which is like `run-parts`.

## run-parts

`run-parts` executes every (executable) file in a directory, in order,
it is used as part of `cron`, and in shell-scripts.  For example you
might have a directory of files to be run every hour and that might
be processed via:

    $ run-parts /etc/cron.hourly/

Optionally `run-parts` can terminate if any of the child-processes exit
with a non-zero exit-code - which is useful for use in shell-scripts.
Unfortunately the CentOS version of `run-parts` doesn't support this ability
and I need it for CI-purposes.


Steve
