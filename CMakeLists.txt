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
        command/check.h)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_link_libraries(Fluent PRIVATE fluent_libc)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)