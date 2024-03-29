
/******************************************************************************
*
*   FILE
*   ----
*   ServerRep.c
*
*   History
*   -------
*   2013-01-15   File created
*
*******************************************************************************
*
*   This file is generated by the 'acplt_builder' command
*
******************************************************************************/


#ifndef OV_COMPILE_LIBRARY_ksbase
#define OV_COMPILE_LIBRARY_ksbase
#endif


#include "ksbase.h"
#include "libov/ov_macros.h"
#include "ksbase_helper.h"


OV_DLLFNCEXPORT OV_STRING ksbase_ServerRep_name_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return pobj->v_name;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_name_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_STRING  value
) {
    return ov_string_setvalue(&pobj->v_name,value);
}

OV_DLLFNCEXPORT OV_STRING* ksbase_ServerRep_protocols_get(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    OV_UINT *pveclen
) {
    *pveclen = pobj->v_protocols.veclen;
    return pobj->v_protocols.value;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_protocols_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_STRING*  value,
    const OV_UINT veclen
) {
    return Ov_SetDynamicVectorValue(&pobj->v_protocols,value,veclen,STRING);
}

OV_DLLFNCEXPORT OV_STRING* ksbase_ServerRep_port_get(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    OV_UINT *pveclen
) {
    *pveclen = pobj->v_port.veclen;
    return pobj->v_port.value;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_port_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_STRING*  value,
    const OV_UINT veclen
) {
    return Ov_SetDynamicVectorValue(&pobj->v_port,value,veclen,STRING);
}

OV_DLLFNCEXPORT OV_UINT ksbase_ServerRep_timetolive_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return pobj->v_timetolive;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_timetolive_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_UINT  value
) {
	OV_TIME exptime;
	OV_TIME_SPAN ttl;

	pobj->v_timetolive = value;

	ttl.secs = value;
	ttl.usecs = 0;

	ov_time_add(&exptime, &(pobj->v_regtime), &ttl);
	ksbase_ServerRep_expirationtime_set(pobj, &exptime);
	return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_INT ksbase_ServerRep_version_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return pobj->v_version;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_version_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_INT  value
) {
    pobj->v_version = value;
    return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_TIME* ksbase_ServerRep_expirationtime_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return &pobj->v_expirationtime;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_expirationtime_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_TIME*  value
) {
    pobj->v_expirationtime = *value;
    return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_TIME* ksbase_ServerRep_regtime_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return &pobj->v_regtime;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_regtime_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_TIME*  value
) {
    OV_TIME exptime;
	pobj->v_regtime = *value;
    exptime = *value;
    exptime.secs += pobj->v_timetolive;
    ksbase_ServerRep_expirationtime_set(pobj, &exptime);
    return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_INT ksbase_ServerRep_state_get(
    OV_INSTPTR_ksbase_ServerRep          pobj
) {
    return pobj->v_state;
}

/**
 * state of the server; 0: offline; 1: online; 2: inactive
 */
OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_state_set(
    OV_INSTPTR_ksbase_ServerRep          pobj,
    const OV_INT  value
) {
    pobj->v_state = value;
    return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_RESULT ksbase_ServerRep_constructor(
	OV_INSTPTR_ov_object 	pobj
) {
	OV_RESULT result;
	OV_INSTPTR_ksbase_ServerRep pthis = Ov_StaticPtrCast(ksbase_ServerRep, pobj);

	/*	baseclass' constructor first	*/
	result = ksbase_ComTask_constructor(pobj);
	if(Ov_Fail(result))
		return result;

	pthis->v_protocols.veclen = 0;
	pthis->v_protocols.value = NULL;
	pthis->v_port.veclen = 0;
	pthis->v_port.value = NULL;

	return OV_ERR_OK;
}

OV_DLLFNCEXPORT void ksbase_ServerRep_startup(
	OV_INSTPTR_ov_object 	pobj
) {
    /*    
    *   local variables
    */

    /* do what the base class does first */
    ov_object_startup(pobj);

    /* do what */


    return;
}

OV_DLLFNCEXPORT void ksbase_ServerRep_shutdown(
	OV_INSTPTR_ov_object 	pobj
) {
    /*    
    *   local variables
    */
    OV_INSTPTR_ksbase_ServerRep pinst = Ov_StaticPtrCast(ksbase_ServerRep, pobj);

    /* do what */
    Ov_SetDynamicVectorLength(&(pinst->v_protocols), 0, STRING);
    Ov_SetDynamicVectorLength(&(pinst->v_port), 0, STRING);
    /* set the object's state to "shut down" */
    ov_object_shutdown(pobj);
    return;
}

/**
* Procedure periodically called by RootComTask
* Checks if timetolife of server has expired.
* When expired server first is set as inactive
* and after another 300 secs object will be deleted.
*/
OV_DLLFNCEXPORT void ksbase_ServerRep_typemethod(
    OV_INSTPTR_ksbase_ComTask          this
) {
	OV_INSTPTR_ksbase_ServerRep pinst = Ov_StaticPtrCast(ksbase_ServerRep, this);
	OV_TIME timenow;

	ov_time_gettime(&timenow);

	if(ov_time_compare(&timenow, &(pinst->v_expirationtime)) == OV_TIMECMP_AFTER) {
		ksbase_ServerRep_state_set(pinst, KSBASE_SERVERREP_STATE_INACTIVE);
	}
	//remove our entry after 5 minutes
	if(timenow.secs<300) // prevent integer underflow
		return;
	timenow.secs -= 300;
	if(ov_time_compare(&timenow, &(pinst->v_expirationtime)) == OV_TIMECMP_AFTER) {
		Ov_DeleteObject(Ov_GetParent(ov_containment, pinst));
	}

	return;
}

