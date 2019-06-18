
/******************************************************************************
*
*   FILE
*   ----
*   ExpressionLogic.c
*
*   History
*   -------
*   2017-06-22   File created
*
*******************************************************************************
*
*   This file is generated by the 'acplt_builder' command
*
******************************************************************************/


#ifndef OV_COMPILE_LIBRARY_propertyValueStatement
#define OV_COMPILE_LIBRARY_propertyValueStatement
#endif


#include "propertyValueStatement.h"
#include "libov/ov_macros.h"


OV_DLLFNCEXPORT OV_ACCESS propertyValueStatement_ExpressionLogic_getaccess(
	OV_INSTPTR_ov_object	pobj,
	const OV_ELEMENT		*pelem,
	const OV_TICKET			*pticket
) {
    /*    
    *   local variables
    */

	switch(pelem->elemtype) {
		case OV_ET_VARIABLE:
			if(pelem->elemunion.pvar->v_offset >= offsetof(OV_INST_ov_object,__classinfo)) {
				if(pelem->elemunion.pvar->v_vartype == OV_VT_CTYPE)
					return OV_AC_NONE;
				else{
					if((pelem->elemunion.pvar->v_varprops & OV_VP_DERIVED)){
						if((pelem->elemunion.pvar->v_varprops & OV_VP_SETACCESSOR)){
							return OV_AC_READWRITE;
						} else {
							return OV_AC_READ;
						}
					} else {
						return OV_AC_READWRITE;
					}
				}
			}
		break;
		default:
		break;
	}

	return ov_object_getaccess(pobj, pelem, pticket);
}

OV_DLLFNCEXPORT OV_RESULT propertyValueStatement_ExpressionLogic_constructor(
	OV_INSTPTR_ov_object 	pobj
) {
    /*    
    *   local variables
    */
    OV_INSTPTR_propertyValueStatement_ExpressionLogic pinst = Ov_StaticPtrCast(propertyValueStatement_ExpressionLogic, pobj);
    OV_RESULT    result;

    /* do what the base class does first */
    result = ov_object_constructor(pobj);
    if(Ov_Fail(result))
         return result;

    /* do what */
    OV_INSTPTR_ov_domain pparent = NULL;
	pparent = Ov_GetParent(ov_containment, pobj);
	if (!Ov_CanCastTo(propertyValueStatement_PropertyValueStatement, pparent) && !Ov_CanCastTo(propertyValueStatement_PropertyValueStatementList, pparent)){
		ov_logfile_error("%s: cannot instantiate - Parent have to be from class propertyValueStatement or propertyValueStatementList", pinst->v_identifier);
		return OV_ERR_ALREADYEXISTS;
	}

	OV_INSTPTR_ov_object pchild = NULL;
	if (Ov_CanCastTo(propertyValueStatement_PropertyValueStatementList, pparent)){
		Ov_ForEachChild(ov_containment, pparent, pchild){
			if (Ov_CanCastTo(propertyValueStatement_ExpressionLogic, pchild) && pchild != pobj){
				ov_logfile_error("%s: cannot instantiate - ExpressionLogic instance already exists", pinst->v_identifier);
				return OV_ERR_ALREADYEXISTS;
			}else if (Ov_CanCastTo(propertyValueStatement_PropertyValueStatement, pchild)){
   				ov_logfile_error("%s: cannot instantiate - at least one instance from class propertyValueStatement already exists", pinst->v_identifier);
				return OV_ERR_ALREADYEXISTS;
   			}
		}
	}

	pchild = NULL;
	if (Ov_CanCastTo(propertyValueStatement_PropertyValueStatement, pparent)){
		Ov_ForEachChild(ov_containment, pparent, pchild){
			if (Ov_CanCastTo(propertyValueStatement_ExpressionLogic, pchild) && pchild != pobj){
				ov_logfile_error("%s: cannot instantiate - ExpressionLogic instance already exists", pinst->v_identifier);
				return OV_ERR_ALREADYEXISTS;
			}
		}
		pparent = Ov_GetParent(ov_containment, pparent);
		pchild = NULL;
		Ov_ForEachChild(ov_containment, pparent, pchild){
			if (Ov_CanCastTo(propertyValueStatement_ExpressionLogic, pchild) && pchild != pobj){
				ov_logfile_error("%s: cannot instantiate - ExpressionLogic instance already exists", pinst->v_identifier);
				return OV_ERR_ALREADYEXISTS;
			}
		}
	}

    return OV_ERR_OK;
}

