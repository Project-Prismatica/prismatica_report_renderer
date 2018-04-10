# MongoDB Connectivity

The templating engine has a filter function called ```mongo``` which will
execute the provided javascript on the specified mongodb instance.

## Example

```bash
The following is from a mongo database:
{% autoescape off %}
{{ "db.myCollection.findOne();" | mongo:"mongo://localhost:27017/myDatabase" }}
{% endautoescape %}
```

## URI Format

The mongoDB URI which the mongo filter function expects is in the following
format:
```
mongo://<hostame:port>/<database name>
```

