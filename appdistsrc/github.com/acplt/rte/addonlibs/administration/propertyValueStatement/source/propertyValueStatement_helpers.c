/*
 * propertyValueStatement_helpers.c
 *
 *  Created on: 07.02.2017
 *      Author: ubuntu
 */

#ifndef OV_COMPILE_LIBRARY_propertyValueStatement
#define OV_COMPILE_LIBRARY_propertyValueStatement
#endif

#include "propertyValueStatement.h"
#include "propertyValueStatement_helpers.h"
#include "libov/ov_database.h"
#include "libov/ov_time.h"
#include "libov/ov_macros.h"

OV_DLLFNCEXPORT PropertyValueStatementList* PropertyValueStatementList_new(){
	PropertyValueStatementList* this = ov_database_malloc(sizeof(PropertyValueStatementList));
	PropertyValueStatementList_init(this);
	return this;
}
OV_DLLFNCEXPORT void PropertyValueStatementList_init(PropertyValueStatementList *this){
	IdentificationType_init(&this->CarrierId);
	this->ExpressionLogic = 0;
	this->ExpressionSemantic = 0;
	IdentificationType_init(&this->ID);
	IdentificationType_init(&this->PropertyId);
	this->Visibility = 0;
	this->View = 0;
	this->Mask = 0;
	this->PvslName = NULL;
	this->pvs = NULL;
	this->pvsSize = 0;
}
OV_DLLFNCEXPORT void PropertyValueStatementList_deleteMembers(PropertyValueStatementList *this){
	IdentificationType_deleteMembers(&this->CarrierId);
	IdentificationType_deleteMembers(&this->ID);
	IdentificationType_deleteMembers(&this->PropertyId);
	for (OV_UINT i = 0; i < this->pvsSize; i++){
		PropertyValueStatement_deleteMembers(&(this->pvs)[i]);
	}
	ov_database_free(this->pvs);
	ov_string_setvalue(&this->PvslName, NULL);
	PropertyValueStatementList_init(this);
}
OV_DLLFNCEXPORT void PropertyValueStatementList_delete(PropertyValueStatementList *this){
	PropertyValueStatementList_deleteMembers(this);
	ov_database_free(this);
}

OV_DLLFNCEXPORT PropertyValueStatement* PropertyValueStatement_new(){
	PropertyValueStatement* this = ov_database_malloc(sizeof(PropertyValueStatement));
	PropertyValueStatement_init(this);
	return this;
}
OV_DLLFNCEXPORT void PropertyValueStatement_init(PropertyValueStatement *this){
	IdentificationType_init(&this->CarrierId);
	this->ExpressionLogic = 0;
	this->ExpressionSemantic = 0;
	IdentificationType_init(&this->ID);
	IdentificationType_init(&this->PropertyId);
	this->Visibility = 0;
	this->View = 0;
	this->Mask = 0;
	this->PvsName = NULL;
	ov_time_gettime(&this->Value.time);
	this->Value.value.vartype = OV_VT_VOID;
	this->Value.value.valueunion.val_string = NULL;
	this->Value.value.valueunion.val_bool_vec.value = NULL;
	this->Value.value.valueunion.val_bool_vec.veclen = 0;
	this->Value.value.valueunion.val_byte_vec.value = NULL;
	this->Value.value.valueunion.val_byte_vec.veclen = 0;
	this->Value.value.valueunion.val_double_vec.value = NULL;
	this->Value.value.valueunion.val_double_vec.veclen = 0;
	this->Value.value.valueunion.val_generic_vec.value = NULL;
	this->Value.value.valueunion.val_generic_vec.veclen = 0;
	this->Value.value.valueunion.val_int_vec.value = NULL;
	this->Value.value.valueunion.val_int_vec.veclen = 0;
	this->Value.value.valueunion.val_single_vec.value = NULL;
	this->Value.value.valueunion.val_single_vec.veclen = 0;
	this->Value.value.valueunion.val_state_vec.value = NULL;
	this->Value.value.valueunion.val_state_vec.veclen = 0;
	this->Value.value.valueunion.val_string_vec.value = NULL;
	this->Value.value.valueunion.val_string_vec.veclen = 0;
	this->Value.value.valueunion.val_struct_vec.value = NULL;
	this->Value.value.valueunion.val_struct_vec.veclen = 0;
	this->Value.value.valueunion.val_time_span_vec.value = NULL;
	this->Value.value.valueunion.val_time_span_vec.veclen = 0;
	this->Value.value.valueunion.val_time_vec.value = NULL;
	this->Value.value.valueunion.val_time_vec.veclen = 0;
	this->Value.value.valueunion.val_uint_vec.value = NULL;
	this->Value.value.valueunion.val_uint_vec.veclen = 0;
}
OV_DLLFNCEXPORT void PropertyValueStatement_deleteMembers(PropertyValueStatement *this){
	IdentificationType_deleteMembers(&this->CarrierId);
	IdentificationType_deleteMembers(&this->ID);
	IdentificationType_deleteMembers(&this->PropertyId);
	ov_string_setvalue(&this->PvsName, NULL);
	Ov_SetAnyValue(&this->Value, NULL);
	PropertyValueStatement_init(this);
}
OV_DLLFNCEXPORT void PropertyValueStatement_delete(PropertyValueStatement *this){
	PropertyValueStatement_deleteMembers(this);
	free(this);
}
