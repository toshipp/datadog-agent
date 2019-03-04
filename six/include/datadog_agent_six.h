// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.
#ifndef DATADOG_AGENT_SIX_H_INCLUDED
#define DATADOG_AGENT_SIX_H_INCLUDED
#include <six_types.h>

#ifdef __cplusplus
extern "C" {
#endif

struct six_s;
typedef struct six_s six_t;

struct six_pyobject_s;
typedef struct six_pyobject_s six_pyobject_t;

// FACTORIES
DATADOG_AGENT_SIX_API six_t *make2();
DATADOG_AGENT_SIX_API six_t *make3();

// API
DATADOG_AGENT_SIX_API void destroy(six_t *);
DATADOG_AGENT_SIX_API int init(six_t *, char *);
DATADOG_AGENT_SIX_API int add_python_path(six_t *, const char *path);
DATADOG_AGENT_SIX_API six_gilstate_t ensure_gil(six_t *);
DATADOG_AGENT_SIX_API void clear_error(six_t *);
DATADOG_AGENT_SIX_API void release_gil(six_t *, six_gilstate_t);
DATADOG_AGENT_SIX_API int get_check(six_t *, const char *name, const char *init_config, const char *instances,
                                    six_pyobject_t **check, char **version);
DATADOG_AGENT_SIX_API const char *run_check(six_t *, six_pyobject_t *check);
DATADOG_AGENT_SIX_API void six_free(six_t *, void *ptr);
DATADOG_AGENT_SIX_API void six_decref(six_t *, six_pyobject_t *);

// CONST API
DATADOG_AGENT_SIX_API int is_initialized(six_t *);
DATADOG_AGENT_SIX_API six_pyobject_t *get_none(const six_t *);
DATADOG_AGENT_SIX_API const char *get_py_version(const six_t *);
DATADOG_AGENT_SIX_API int run_simple_string(const six_t *, const char *code);
DATADOG_AGENT_SIX_API int has_error(const six_t *);
DATADOG_AGENT_SIX_API const char *get_error(const six_t *);

// AGGREGATOR API
DATADOG_AGENT_SIX_API void set_submit_metric_cb(six_t *, cb_submit_metric_t);
DATADOG_AGENT_SIX_API void set_submit_service_check_cb(six_t *, cb_submit_service_check_t);
DATADOG_AGENT_SIX_API void set_submit_event_cb(six_t *, cb_submit_event_t);

// DATADOG_AGENT API
DATADOG_AGENT_SIX_API void set_get_version_cb(six_t *, cb_get_version_t);
DATADOG_AGENT_SIX_API void set_get_config_cb(six_t *, cb_get_config_t);
DATADOG_AGENT_SIX_API void set_headers_cb(six_t *, cb_headers_t);
DATADOG_AGENT_SIX_API void set_get_hostname_cb(six_t *, cb_get_hostname_t);
DATADOG_AGENT_SIX_API void set_get_clustername_cb(six_t *, cb_get_clustername_t);
DATADOG_AGENT_SIX_API void set_log_cb(six_t *, cb_log_t);

#ifdef __cplusplus
}
#endif
#endif