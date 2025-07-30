include(FetchContent)

cmake_minimum_required(VERSION 3.31)
project(Fluent)

FetchContent_Declare(
        fluent_libc
        GIT_REPOSITORY https://github.com/rodrigoo-r/fluent_libc
        GIT_TAG        master
)

FetchContent_Declare(
        unordered_dense
        GIT_REPOSITORY https://github.com/martinus/unordered_dense
        GIT_TAG        main
)

FetchContent_Declare(
        xxhash
        GIT_REPOSITORY https://github.com/Cyan4973/xxHash.git
        GIT_TAG        v0.8.3
)

FetchContent_MakeAvailable(xxhash)
FetchContent_MakeAvailable(fluent_libc)
FetchContent_MakeAvailable(unordered_dense)
set(CMAKE_CXX_STANDARD 20)

add_executable(Fluent
        main.cpp
        container/owned_string.h
        container/vector.h
        container/external_string.h
        except/exception.h
        cli/args.h
        cli/app.h
        cli/value.h
)

target_include_directories(Fluent PRIVATE ${CMAKE_BINARY_DIR}/_deps/fluent_libc-src/include)
target_include_directories(Fluent PRIVATE .)
target_link_libraries(Fluent PRIVATE fluent_libc)
target_link_libraries(Fluent PRIVATE unordered_dense::unordered_dense)
target_link_libraries(Fluent PRIVATE xxhash)
add_definitions(-DDEBUG -DFLUENT_LIBC_RELEASE=1)
set_target_properties(Fluent PROPERTIES OUTPUT_NAME "fluent")

if (NOT WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -march=native -funroll-loops -flto")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O3 -march=native -funroll-loops -flto")
endif ()