include(FetchContent)

cmake_minimum_required(VERSION 3.31)
project(Zelix)

FetchContent_Declare(
        zelix_stl
        GIT_REPOSITORY https://github.com/zelix-lang/STL.git
        GIT_TAG        master
)

FetchContent_Declare(
        fluent_libc
        GIT_REPOSITORY https://github.com/rodrigoo-r/fluent_libc
        GIT_TAG        master
)

FetchContent_Declare(
        zelix_cli
        GIT_REPOSITORY https://github.com/zelix-lang/CLI.git
        GIT_TAG        master
)

FetchContent_MakeAvailable(fluent_libc)
FetchContent_MakeAvailable(zelix_cli)
FetchContent_MakeAvailable(unordered_dense)
FetchContent_MakeAvailable(zelix_stl)
set(CMAKE_CXX_STANDARD 20)

add_executable(Zelix
        main.cpp
        lexer/lexer.h
        lexer/token.h
        lexer/lexer.cpp
        util/read_file.h
        parser/ast.h
        parser/parser.h
        parser/parser.cpp
        parser/rule/import.h
        parser/expect.h
        parser/rule/function.h
        parser/rule/type/type.h
        memory/allocator.h
        parser/rule/block/block.h
        parser/rule/block/block.cpp
        parser/rule/expr/expr.h
        parser/rule/extractor.h
        parser/rule/expr/queue.h
        parser/rule/expr/likely.h
        parser/rule/args.h
        parser/rule/call.h
        parser/rule/prop.h
        parser/rule/signed.h
        command/compile.h
        time/timed_task.h
        time/timed_task.cpp
        parser/rule/declaration.h
        parser/rule/if.h
        parser/rule/for.h
        parser/rule/mod.h
        analyzer/analyzer.h
        parser/rule/derive.h
        parser/rule/assignment.h
        list/list.h
        file_code/file_code.h
        file_code/function.h
        file_code/mod.h
        file_code/type.h
        parser/rule/type/type.cpp
        file_code/declaration.h
        file_code/converter/converter.h
        file_code/converter/function.h
        file_code/converter/type/type.h
        file_code/converter/type/type.cpp
        file_code/converter/converter.cpp
        file_code/converter/import/import.h
        global/constants.h
        util/dirname.h
        util/absolute.h
        file_code/converter/mod.h
)

target_include_directories(Zelix PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_include_directories(Zelix PRIVATE .)
target_link_libraries(Zelix PRIVATE fluent_libc)
target_link_libraries(Zelix PRIVATE zelix::stl)
target_link_libraries(Zelix PRIVATE zelix::cli)
target_link_libraries(Zelix PRIVATE unordered_dense::unordered_dense)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)
set_target_properties(Zelix PROPERTIES OUTPUT_NAME "zelix")

if (NOT WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -march=native -funroll-loops -flto")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O3 -march=native -funroll-loops -flto")
endif ()