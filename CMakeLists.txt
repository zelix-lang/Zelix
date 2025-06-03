include(FetchContent)

FetchContent_Declare(
        fluent_libc
        GIT_REPOSITORY https://github.com/rodrigoo-r/fluent_libc
        GIT_TAG        master
)

FetchContent_MakeAvailable(fluent_libc)
cmake_minimum_required(VERSION 3.30)
project(Fluent C)

set(CMAKE_C_STANDARD 11)

add_executable(Fluent
    main.c
        command/check.h
        lexer/lexer.h
        token/token.h
        token/token_map.h
        lexer/stream.h
        lexer/error.h
        file/file_reader.h)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_link_libraries(Fluent PRIVATE fluent_libc)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)

# Add aggressive optimization flags for C++ and C compilers
if (NOT WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -march=native -funroll-loops -flto")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O3 -march=native -funroll-loops -flto")
endif ()

