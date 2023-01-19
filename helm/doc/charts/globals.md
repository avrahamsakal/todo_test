# Configure Charts using Globals

To reduce configuration duplication when installing our wrapper Helm chart, several
configuration settings are available to be set in the `global` section of `values.yml`.
These global settings are used across several charts, while all other settings are scoped
within their chart. See the [Helm documentation on globals](https://docs.helm.sh/developing_charts/#global-values)
for more information on how the global variables work.

- [MYSQL](#configure-mysql-settings)

## Configure MYSQL settings

The GitLab global PostgreSQL settings are located under the `global.mysql` key.

```YAML
global:
  psql:
    host: db.example.com
    port: 5432
    password:
      secret: gitlab-postgres
      key: psql-password
```

If you want to connect Gitlab with a PostgreSQL database over mutual TLS, create a secret
containing the client key, client certificate and server certificate authority as different
secret keys. Then describe the secret's structure using the `global.psql.ssl` mapping.

```YAML
global:
  psql:
    host: db.example.com
    # ... further settings like in the previous example ...
    ssl:
      secret: db-example-ssl-secrets # Name of the secret
      clientKey: key.pem             # Secret key of the certificate's key
      clientCertificate: cert.pem    # Secret key storing the certificate
      serverCA: server-ca.pem        # Secret key containing the CA for the database server
```
