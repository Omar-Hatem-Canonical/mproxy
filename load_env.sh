export MPROXY_MQTT_WITHOUT_TLS_ADDRESS=":1884"
export MPROXY_MQTT_WITHOUT_TLS_TARGET="localhost:1883"

export MPROXY_MQTT_WITH_TLS_ADDRESS=":8883"
export MPROXY_MQTT_WITH_TLS_TARGET="localhost:1883"
export MPROXY_MQTT_WITH_TLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_MQTT_WITH_TLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_MQTT_WITH_TLS_SERVER_CA_FILE="ssl/certs/ca.crt"
 
export MPROXY_MQTT_WITH_MTLS_ADDRESS=":8884"
export MPROXY_MQTT_WITH_MTLS_TARGET="localhost:1883"
export MPROXY_MQTT_WITH_MTLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_MQTT_WITH_MTLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_MQTT_WITH_MTLS_SERVER_CA_FILE="ssl/certs/ca.crt"
export MPROXY_MQTT_WITH_MTLS_CLIENT_CA_FILE="ssl/certs/ca.crt"
export MPROXY_MQTT_WITH_MTLS_CERT_VERIFICATION_METHODS="ocsp"
export MPROXY_MQTT_WITH_MTLS_OCSP_RESPONDER_URL="http://localhost:8080/ocsp"
 
export MPROXY_MQTT_WS_WITHOUT_TLS_ADDRESS=":8083"
export MPROXY_MQTT_WS_WITHOUT_TLS_TARGET="ws://localhost:8000/"
 
export MPROXY_MQTT_WS_WITH_TLS_ADDRESS=":8084"
export MPROXY_MQTT_WS_WITH_TLS_TARGET="ws://localhost:8000/"
export MPROXY_MQTT_WS_WITH_TLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_MQTT_WS_WITH_TLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_MQTT_WS_WITH_TLS_SERVER_CA_FILE="ssl/certs/ca.crt"
 
export MPROXY_MQTT_WS_WITH_MTLS_ADDRESS=":8085"
export MPROXY_MQTT_WS_WITH_MTLS_PATH_PREFIX="/mqtt"
export MPROXY_MQTT_WS_WITH_MTLS_TARGET="ws://localhost:8000/"
export MPROXY_MQTT_WS_WITH_MTLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_MQTT_WS_WITH_MTLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_MQTT_WS_WITH_MTLS_SERVER_CA_FILE="ssl/certs/ca.crt"
export MPROXY_MQTT_WS_WITH_MTLS_CLIENT_CA_FILE="ssl/certs/ca.crt"
export MPROXY_MQTT_WS_WITH_MTLS_CERT_VERIFICATION_METHODS="ocsp"
export MPROXY_MQTT_WS_WITH_MTLS_OCSP_RESPONDER_URL="http://localhost:8080/ocsp"
 
export MPROXY_HTTP_WITHOUT_TLS_ADDRESS=":8086"
export MPROXY_HTTP_WITHOUT_TLS_PATH_PREFIX="/messages"
export MPROXY_HTTP_WITHOUT_TLS_TARGET="http://localhost:8888/"
 
export MPROXY_HTTP_WITH_TLS_ADDRESS=":8087"
export MPROXY_HTTP_WITH_TLS_PATH_PREFIX="/messages"
export MPROXY_HTTP_WITH_TLS_TARGET="http://localhost:8888/"
export MPROXY_HTTP_WITH_TLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_HTTP_WITH_TLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_HTTP_WITH_TLS_SERVER_CA_FILE="ssl/certs/ca.crt"
 
export MPROXY_HTTP_WITH_MTLS_ADDRESS=":8088"
export MPROXY_HTTP_WITH_MTLS_PATH_PREFIX="/messages"
export MPROXY_HTTP_WITH_MTLS_TARGET="http://localhost:8888/"
export MPROXY_HTTP_WITH_MTLS_CERT_FILE="ssl/certs/server.crt"
export MPROXY_HTTP_WITH_MTLS_KEY_FILE="ssl/certs/server.key"
export MPROXY_HTTP_WITH_MTLS_SERVER_CA_FILE="ssl/certs/ca.crt"
export MPROXY_HTTP_WITH_MTLS_CLIENT_CA_FILE="ssl/certs/ca.crt"
export MPROXY_HTTP_WITH_MTLS_CERT_VERIFICATION_METHODS="ocsp"
export MPROXY_HTTP_WITH_MTLS_OCSP_RESPONDER_URL="http://localhost:8080/ocsp"
