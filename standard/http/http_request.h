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

#ifndef HTTP_REQUEST_H
#define HTTP_REQUEST_H

#include <string>
#include <map>
#include "message_encoding.cpp"

class HttpRequest {
    private:
        MessageEncoding message_encoding;
        std::string body;
        std::string origin;
        std::string path;
        std::string method;
        std::map<std::string, std::string> headers;

    public:
        HttpRequest(
            MessageEncoding message_encoding,
            std::string body,
            std::string origin,
            std::string path,
            std::string method,
            std::map<std::string, std::string> headers
        );

        MessageEncoding get_message_encoding();
        std::string get_body();
        std::string get_origin();
        std::string get_path();
        std::string get_method();
        std::map<std::string, std::string> get_headers();
};

#endif
