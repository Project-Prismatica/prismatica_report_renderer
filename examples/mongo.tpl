The following is from a mongo database:
{% autoescape off %}
{{ "db.myCollection.findOne();" | mongo:"mongo://localhost:27017/myDatabase" }}
{% endautoescape %}
