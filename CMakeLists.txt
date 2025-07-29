include(FetchContent)

cmake_minimum_required(VERSION 3.31)
project(Fluent)

FetchContent_Declare(
        fluent_libc
        GIT_REPOSITORY https://github.com/rodrigoo-r/fluent_libc
        GIT_TAG        master
)

FetchContent_MakeAvailable(fluent_libc)
set(CMAKE_CXX_STANDARD 20)

add_executable(Fluent
        main.cpp
        container/owned_string.h
)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_link_libraries(Fluent PRIVATE fluent_libc)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)