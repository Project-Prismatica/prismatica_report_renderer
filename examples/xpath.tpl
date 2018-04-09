{% with inputXml="<note><togroup><to>Tove</to><to>Alice</to></togroup><from>Jani</from><heading>Reminder</heading><body>forget me not this weekend!</body></note>" %}
Recipients:
{% for recipient in inputXml | xpath:"/note/togroup/*" %}* {{ recipient }}
{% endfor %}
From: {{ inputXml | xpath:"/note/from" }}
{% autoescape off %}
Body: {{ inputXml | xpath:"/note/body" }}
{% endautoescape %}
{% endwith %}