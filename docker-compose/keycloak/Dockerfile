FROM quay.io/keycloak/keycloak:26.3.2

# Copy custom Portcall themes
COPY themes/portcall /opt/keycloak/themes/portcall
COPY themes/portcall-admin /opt/keycloak/themes/portcall-admin

# Copy realm configurations
COPY realms/dev-realm.json /opt/keycloak/data/import/dev-realm.json
COPY realms/admin-realm.json /opt/keycloak/data/import/admin-realm.json

# Build optimized Keycloak image with database and email template settings
RUN /opt/keycloak/bin/kc.sh build --db=postgres --spi-email-template-provider=freemarker --spi-email-template-freemarker-enabled=true
