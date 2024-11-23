#ifndef HTTP_SERVER_H
#define HTTP_SERVER_H

#include "http_response.h"
#include "http_request.h"
#include "http_server.h"
#include "../lang/result.h"

Result<bool> create_http_server(int port, HttpResponse (*callback)(HttpRequest));

#endif