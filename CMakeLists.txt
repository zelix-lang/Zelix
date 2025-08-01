include(FetchContent)

cmake_minimum_required(VERSION 3.31)
project(Fluent)

FetchContent_Declare(
        fluent_stl
        GIT_REPOSITORY https://github.com/fluent-lang/STL.git
        GIT_TAG        master
)

FetchContent_Declare(
        fluent_libc
        GIT_REPOSITORY https://github.com/rodrigoo-r/fluent_libc
        GIT_TAG        master
)

FetchContent_Declare(
        fluent_cli
        GIT_REPOSITORY https://github.com/fluent-lang/CLI.git
        GIT_TAG        master
)

FetchContent_MakeAvailable(fluent_libc)
FetchContent_MakeAvailable(fluent_cli)
FetchContent_MakeAvailable(unordered_dense)
FetchContent_MakeAvailable(fluent_stl)
set(CMAKE_CXX_STANDARD 20)

add_executable(Fluent
        main.cpp
        lexer/lexer.h
        lexer/token.h
        lexer/lexer.cpp
        util/read_file.h
        parser/ast.h
        parser/parser.h
        parser/parser.cpp)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_include_directories(Fluent PRIVATE .)
target_link_libraries(Fluent PRIVATE fluent_libc)
target_link_libraries(Fluent PRIVATE fluent::stl)
target_link_libraries(Fluent PRIVATE fluent::cli)
target_link_libraries(Fluent PRIVATE unordered_dense::unordered_dense)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)
set_target_properties(Fluent PROPERTIES OUTPUT_NAME "fluent")

if (NOT WIN32)
#    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -march=native -funroll-loops -flto")
#    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O3 -march=native -funroll-loops -flto")
endif ()