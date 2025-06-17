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
        file/file_reader.h
        logger/logger.h
        state/state.h
        message/generator.h
        message/lexer_error_converter.h
        ast/ast.h
        ast/rule.h
        parser/parser.h
        parser/error.h
        parser/stream.h
        parser/rule/import.h
        parser/extractor.h
        parser/rule/function.h
        type/type.h
        parser/rule/type.h
        parser/rule/block.h
        parser/rule/expression.h
        parser/queue/expression.h)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_link_libraries(Fluent PRIVATE fluent_libc m)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)

# Add aggressive optimization flags for C++ and C compilers
if (NOT WIN32)
#    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -march=native -funroll-loops -flto")
#    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O3 -march=native -funroll-loops -flto")
endif ()

# Jemalloc (fluent_libc-defined)
if(NOT WIN32) # Guard against Windows, falls to libc on Windows
    find_package(PkgConfig REQUIRED)
    pkg_check_modules(JEMALLOC QUIET jemalloc)
    if(JEMALLOC_FOUND)
        target_compile_definitions(Fluent PUBLIC DEBUG HAVE_JEMALLOC=1)
        message(STATUS "jemalloc found! Adding it to the link libraries.")
        target_include_directories(Fluent PRIVATE ${JEMALLOC_INCLUDE_DIRS})
        target_link_libraries(Fluent PRIVATE ${JEMALLOC_LIBRARIES})
    else()
        message(STATUS "jemalloc NOT found, skipping.")
    endif()
endif()