// Copyright (c) 2014 The cef2go authors. All rights reserved.
// License: BSD 3-clause.
// Website: https://github.com/fromkeith/cef2go
package cef2go

/*
#include "cefBase.h"
#include "include/capi/cef_scheme_capi.h"

extern void intialize_cef_scheme_handler_factory(struct _cef_scheme_handler_factory_t * factory);
extern int cef_scheme_handler_register(char * schemeName, char * domainName, struct _cef_scheme_handler_factory_t* factory);
extern int whitelist_cef_add_cross_origin_whitelist_entry(char* sourceOrigin, char* targetProtocol, char* targetDomain, int allow_target_subdomains);
*/
import "C"

import (
    "unsafe"
    //"fmt"
)


type CefSchemeHandlerFactory struct {
    CStruct         *C.struct__cef_scheme_handler_factory_t // memory manage?
    ResourceHandler     CefResourceHandlerT
}

var (
    schemeHandlerMap = make(map[*C.struct__cef_scheme_handler_factory_t]CefSchemeHandlerFactory)
)


//export go_CreateSchemeHandler
func go_CreateSchemeHandler(
        self *C.struct__cef_scheme_handler_factory_t,
        browser *C.struct__cef_browser_t,
        frame *C.struct__cef_frame_t,
        scheme_name *C.cef_string_utf8_t,
        request *C.struct__cef_request_t) *C.struct__cef_resource_handler_t {

    defer C.cef_string_userfree_utf8_free(scheme_name)

    if handler, ok := schemeHandlerMap[self]; ok {
        return handler.ResourceHandler.CStruct;
    }
    return nil
}

func RegisterCustomScheme(schemeName, domainName string, resHandler ResourceHandler) (int, CefSchemeHandlerFactory) {
    var handler CefSchemeHandlerFactory

    handler.ResourceHandler = createResourceHandler(resHandler)

    handler.CStruct = (*C.struct__cef_scheme_handler_factory_t)(
            C.calloc(1, C.sizeof_struct__cef_scheme_handler_factory_t))
    C.intialize_cef_scheme_handler_factory(handler.CStruct)

    schemeNameCs := C.CString(schemeName)
    defer C.free(unsafe.Pointer(schemeNameCs))

    domainNameCs := C.CString(domainName)
    defer C.free(unsafe.Pointer(domainNameCs))

    retCode := C.cef_scheme_handler_register(schemeNameCs, domainNameCs, handler.CStruct)

    schemeHandlerMap[handler.CStruct] = handler

    return int(retCode), handler
}

func AddCrossOriginWhitelistEntry(sourceOrigin, targetProtocol, targetDomain string, allowSubdomains bool) bool {

    sourceOriginCs := C.CString(sourceOrigin)
    defer C.free(unsafe.Pointer(sourceOriginCs))
    targetProtocolCs := C.CString(targetProtocol)
    defer C.free(unsafe.Pointer(targetProtocolCs))
    targetDomainCs := C.CString(targetDomain)
    defer C.free(unsafe.Pointer(targetDomainCs))

    allowDomainsInt := 1
    if !allowSubdomains {
        allowDomainsInt = 0
    }

    return C.whitelist_cef_add_cross_origin_whitelist_entry(sourceOriginCs, targetProtocolCs, targetDomainCs, C.int(allowDomainsInt)) == C.int(1)

}