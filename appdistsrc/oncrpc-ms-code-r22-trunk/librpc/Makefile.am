## Process this file with automake to produce Makefile.in

lib_LTLIBRARIES = liboncrpc.la

AM_CFLAGS = -D_X86_ -DONCRPCDLL
INCLUDES = -I$(top_srcdir)

# MinGW requires -no-undefined if a DLL is to be built.
liboncrpc_la_LDFLAGS = -no-undefined -version-info @SHARED_VERSION_INFO@ @SHLIB_VERSION_ARG@ -lws2_32
liboncrpc_la_SOURCES = \
    clnt_raw.c \
    xdr.c \
    clnt_tcp.c \
    clnt_udp.c \
    pmap_rmt.c \
    rpc_prot.c \
    svc_autu.c \
    svc_auth.c \
    svc_raw.c \
    svc_run.c \
    svc_tcp.c \
    svc_udp.c \
    xdr_mem.c \
    xdr_rec.c \
    auth_non.c \
    auth_uni.c \
    authunix.c \
    bindresv.c \
    clnt_gen.c \
    clnt_per.c \
    clnt_sim.c \
    get_myad.c \
    getrpcen.c \
    getrpcpo.c \
    pmap_cln.c \
    pmap_get.c \
    pmap_gma.c \
    pmap_pro.c \
    pmap_pr.c \
    rpc_call.c \
    rpc_comm.c \
    svc_simp.c \
    xdr_arra.c \
    xdr_floa.c \
    xdr_refe.c \
    xdr_stdi.c \
    svc.c \
    bcopy.c \
    nt.c
