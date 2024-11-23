/*
    These files are part of the Surf's standard library.
    They're bundled with the main program by the compiler
    which is then converted to machine code.

    -----
    License notice:

    This code is released under the GNU GPL v3 license.
    The code is provided as is without any warranty
    Copyright (c) 2024 Rodrigo R. & all Surf contributors
*/

#include "http_response.h"
#include <map>

HttpResponse::HttpResponse(
    MessageEncoding message_encoding,
    std::string message,
    std::string status_detail,
    std::map<std::string, std::string> headers,
    int status_code
) :
    message_encoding(message_encoding),
    message(message),
    status_detail(status_detail),
    headers(headers),
    status_code(status_code) {}

std::string HttpResponse::get_message() {
    return message;
}

MessageEncoding HttpResponse::get_message_encoding() {
    return message_encoding;
}

int HttpResponse::get_status_code() {
    return status_code;
}

std::string HttpResponse::get_status_detail() {
    return status_detail;
}

std::map<std::string, std::string> HttpResponse::get_headers() {
    return headers;
}