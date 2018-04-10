# XML Support

The report renderer exposes an ```xpath``` filter function for use in templates.
It is designed to run the specified xpath query on the provided XML string. In
cases where more than one result is returned, one can use a for loop to iterate
over results.

Example may be found under ```examples/xpath.tpl```.

# Quirks

## XPATH attributes

Currently, with an xpath attribute query (like ```/report/finding/@ctime``),
it will only return the contents of the XML node with the attribute ```ctime```.
This currently looks like an idiosyncracy with the
(xpath library)[https://github.com/antchfx/xpath] currently in employ.
