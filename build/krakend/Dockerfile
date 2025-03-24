FROM devopsfaith/krakend:2.9.3

COPY ./krakend.tmpl /etc/krakend/krakend.tmpl
COPY ./config/ /etc/krakend/config/
COPY ./plugin/deprecated-header.so /opt/krakend/plugins/deprecated-header.so

ENV FC_ENABLE=1
ENV FC_PARTIALS=/etc/krakend/config/partials
ENV FC_SETTINGS=/etc/krakend/config/settings
ENV FC_TEMPLATES=/etc/krakend/config/templates

WORKDIR /etc/krakend

CMD ["run","-dc","krakend.tmpl"]