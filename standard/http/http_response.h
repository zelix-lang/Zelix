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

#ifndef HTTP_RESPONSE_H
#define HTTP_RESPONSE_H

#include <string>
#include <map>
#include "message_encoding.cpp"

class HttpResponse {
    private:
        MessageEncoding message_encoding;
        std::string message;
        std::string status_detail;
        std::map<std::string, std::string> headers;
        int status_code;
    
    public:
        HttpResponse(
            MessageEncoding message_encoding,
            std::string message,
            std::string status_detail,
            std::map<std::string, std::string> headers,
            int status_code
        );

        std::string get_message();
        MessageEncoding get_message_encoding();
        int get_status_code();
        std::string get_status_detail();
        std::map<std::string, std::string> get_headers();
};

#endif
