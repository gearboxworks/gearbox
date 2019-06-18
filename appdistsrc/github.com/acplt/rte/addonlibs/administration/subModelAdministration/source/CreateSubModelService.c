
/******************************************************************************
*
*   FILE
*   ----
*   CreateSubModelService.c
*
*   History
*   -------
*   2017-10-10   File created
*
*******************************************************************************
*
*   This file is generated by the 'acplt_builder' command
*
******************************************************************************/


#ifndef OV_COMPILE_LIBRARY_subModelAdministration
#define OV_COMPILE_LIBRARY_subModelAdministration
#endif


#include "subModelAdministration.h"
#include "libov/ov_macros.h"
#include "helper.h"


OV_DLLFNCEXPORT OV_RESULT subModelAdministration_CreateSubModelService_CallMethod(      
  OV_INSTPTR_services_Service pobj,       
  OV_UINT numberofInputArgs,       
  const void **packedInputArgList,       
  OV_UINT numberofOutputArgs,      
  void **packedOutputArgList,
  OV_UINT *typeArray       
) {
    /*    
    *   local variables
    */
	OV_UINT result = 0;
	OV_INSTPTR_ov_object pParent = NULL;
	OV_INSTPTR_openaas_SubModel pSubModel = NULL;
	OV_STRING status = NULL;

	packedOutputArgList[0] = ov_database_malloc(sizeof(OV_STRING));
	*(OV_STRING*)packedOutputArgList[0] = NULL;
	typeArray[0] = OV_VT_STRING;


	if (*(OV_UINT*)packedInputArgList[1] != 0){
		ov_string_setvalue(&status, "Only Uri as identifier are implemented");
		goto FINALIZE;
	}

	pParent = ov_path_getobjectpointer(*(OV_STRING*)(packedInputArgList[0]), 2);
	if (pParent == NULL){
		ov_string_setvalue(&status, "Find no Object for parrentId");
		goto FINALIZE;
	}

	result = (OV_UINT)checkForSameAAS(Ov_PtrUpCast(ov_object, pobj), pParent);
	if (result != OV_ERR_OK){
		ov_string_setvalue(&status, "Parent is not in the same AAS as the method");
		goto FINALIZE;
	}

	result = Ov_CreateObject(openaas_SubModel, pSubModel, Ov_StaticPtrCast(ov_domain, pParent), *(OV_STRING*)(packedInputArgList[4]));
	if(Ov_Fail(result)){
		ov_logfile_error("Fatal: could not create SubModel object - reason: %s", ov_result_getresulttext(result));
		ov_string_setvalue(&status, ov_result_getresulttext(result));
		goto FINALIZE;
	}

	ov_string_setvalue(&pSubModel->p_ModelId.v_IdSpec, *(OV_STRING*)(packedInputArgList[2]));
	pSubModel->p_ModelId.v_IdType = *(OV_UINT*)packedInputArgList[3];
	pSubModel->v_Revision = *(OV_UINT*)packedInputArgList[5];
	pSubModel->v_Version = *(OV_UINT*)packedInputArgList[6];

	ov_string_setvalue(&status, ov_result_getresulttext(result));
	FINALIZE:

	*(OV_STRING*)packedOutputArgList[0] = ov_database_malloc(ov_string_getlength(status)+1);
	ov_string_setvalue((OV_STRING*)packedOutputArgList[0],status);
	ov_string_setvalue(&status,NULL);

    return OV_ERR_OK;
}

