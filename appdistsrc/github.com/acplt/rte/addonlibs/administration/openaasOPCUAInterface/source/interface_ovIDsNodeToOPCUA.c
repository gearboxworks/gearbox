/******************************************************************************
 *
 *   FILE
 *   ----
 *   nodeStoreFunctions.c
 *
 *   History
 *   -------
 *   2014-10-21   File created
 *
 *******************************************************************************
 *
 *   This file is generated by the 'acplt_builder' command
 *
 ******************************************************************************/

#ifndef OV_COMPILE_LIBRARY_openaasOPCUAInterface
#define OV_COMPILE_LIBRARY_openaasOPCUAInterface
#endif

#include "openaasOPCUAInterface.h"
#include "libov/ov_macros.h"
#include "ksbase.h"
#include "opcua.h"
#include "opcua_helpers.h"
#include "NoneTicketAuthenticator.h"
#include "libov/ov_path.h"
#include "libov/ov_memstack.h"
#include "ks_logfile.h"
#include "nodeset_openaas.h"
#include "nodeset_identification.h"
#include "nodeset_lifeCycleEntry.h"

extern OV_INSTPTR_openaasOPCUAInterface_interface pinterface;


OV_DLLFNCEXPORT UA_StatusCode openaasOPCUAInterface_interface_ovIDsNodeToOPCUA(
		void *handle, const UA_NodeId *nodeId, UA_Node** opcuaNode) {
	UA_Node 				*newNode = NULL;
	UA_StatusCode 			result = UA_STATUSCODE_GOOD;
	OV_PATH 				path;
	OV_INSTPTR_ov_object	pobj = NULL;
	OV_TICKET 				*pTicket = NULL;
	OV_VTBLPTR_ov_object	pVtblObj = NULL;
	OV_ACCESS				access;
	UA_NodeClass 			nodeClass;
	OV_ELEMENT				element;

	if (pinterface == NULL)
		return UA_STATUSCODE_BADOUTOFSERVICE;

	ov_memstack_lock();
	result = opcua_nodeStoreFunctions_resolveNodeIdToPath(*nodeId, &path);
	if(result != UA_STATUSCODE_GOOD){
		ov_memstack_unlock();
		return result;
	}
	element = path.elements[path.size-1];
	ov_memstack_unlock();
	result = opcua_nodeStoreFunctions_getVtblPointerAndCheckAccess(&(element), pTicket, &pobj, &pVtblObj, &access);
	if(result != UA_STATUSCODE_GOOD){
		return result;
	}


	nodeClass = UA_NODECLASS_VARIABLE;
	newNode = (UA_Node*)UA_calloc(1, sizeof(UA_VariableNode));


	// Basic Attribute
	// BrowseName
	UA_QualifiedName qName;
	qName.name = UA_String_fromChars(pobj->v_identifier);
	qName.namespaceIndex = nodeId->namespaceIndex;
	newNode->browseName = qName;

	// Description
	OV_STRING tempString = pVtblObj->m_getcomment(pobj, &element);
	UA_LocalizedText lText;
	UA_LocalizedText_init(&lText);
	lText.locale = UA_String_fromChars("en");
	if(tempString){
		lText.text = UA_String_fromChars(tempString);
	} else {
		lText.text = UA_String_fromChars("");
	}
	UA_LocalizedText_copy(&lText,&newNode->description);
	UA_LocalizedText_deleteMembers(&lText);

	// DisplayName
	UA_LocalizedText displayName;
	UA_LocalizedText_init(&displayName);
	displayName.locale = UA_String_fromChars("en");
	displayName.text = UA_String_fromChars(pobj->v_identifier);
	UA_LocalizedText_copy(&displayName, &newNode->displayName);
	UA_LocalizedText_deleteMembers(&displayName);

	// NodeId
	UA_NodeId_copy(nodeId, &newNode->nodeId);

	// NodeClass
	newNode->nodeClass 	= nodeClass;

	// WriteMask
	UA_UInt32 writeMask = 0;
	if(element.elemtype != OV_ET_VARIABLE){
		if(access & OV_AC_WRITE){
			writeMask |= (1<<2);	//	BrowseName
			writeMask |= (1<<6);	//	DisplayName
		}
		if(access & OV_AC_RENAMEABLE){
			writeMask |= (1<<14);	//	NodeId
		}
	}
	newNode->writeMask 	= writeMask;

	// Variable specific attributes
	// arrayDemensions
	((UA_VariableNode*)newNode)->arrayDimensionsSize = 0;
	((UA_VariableNode*)newNode)->arrayDimensions = NULL; // UA_Array_new(((UA_VariableNode*)newNode)->arrayDimensionsSize, &UA_TYPES[UA_TYPES_INT32]);	/*	scalar or one dimension	*/

	// valuerank
	((UA_VariableNode*)newNode)->valueRank = 1;	/*	one dimension	*/


	// value
	OV_ELEMENT tmpPart;
	tmpPart.elemtype = OV_ET_NONE;
	OV_ELEMENT tmpParrent;
	tmpParrent.pobj = pobj;
	tmpParrent.elemtype = OV_ET_OBJECT;
	UA_Identification tmpidentification;
	UA_Identification_init(&tmpidentification);
	do {
		ov_element_getnextpart(&tmpParrent, &tmpPart, OV_ET_VARIABLE);
		if (tmpPart.elemtype == OV_ET_NONE)
			break;
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "IdSpec") == OV_STRCMP_EQUAL){
			if (*(OV_STRING*)tmpPart.pvalue != NULL)
				tmpidentification.idSpec = UA_String_fromChars(*(OV_STRING*)tmpPart.pvalue);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "IdType") == OV_STRCMP_EQUAL){
			tmpidentification.idType = *(UA_UInt32*)tmpPart.pvalue;
			continue;
		}
	} while(TRUE);


	((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->type = &UA_IDENTIFICATION[UA_IDENTIFICATION_IDENTIFICATION];
	((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data = UA_Identification_new();
	if (!((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data){
		result = UA_STATUSCODE_BADOUTOFMEMORY;
		return result;
	}
	((UA_VariableNode*)newNode)->value.data.value.hasValue = TRUE;
	((UA_VariableNode*)newNode)->valueSource = UA_VALUESOURCE_DATA;
	UA_Identification_copy(&tmpidentification, ((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data);
	UA_Identification_deleteMembers(&tmpidentification);
	// dataType
	((UA_VariableNode*)newNode)->dataType = UA_NODEID_NUMERIC(pinterface->v_modelnamespaceIndexIdentification, UA_NSIDENTIFICATIONID_IDENTIFICATION);


	// accessLevel
	UA_Byte accessLevel = 0;
	if(access & OV_AC_READ){
		accessLevel |= (1<<0);
	}
	if(access & OV_AC_WRITE){
		accessLevel |= (1<<1);
	}
	((UA_VariableNode*)newNode)->accessLevel = accessLevel;
	// minimumSamplingInterval
	((UA_VariableNode*)newNode)->minimumSamplingInterval = -1;
	// historizing
	((UA_VariableNode*)newNode)->historizing = UA_FALSE;

	// References
	addReference(newNode);

	OV_UINT len = 0;
	OV_STRING *plist = NULL;
	OV_STRING tmpString = NULL;
	copyOPCUAStringToOV(nodeId->identifier.string, &tmpString);
	plist = ov_string_split(tmpString, ".", &len);
	ov_string_setvalue(&tmpString, plist[0]);
	for (OV_UINT i = 1; i < len-1; i++){
		ov_string_append(&tmpString, ".");
		ov_string_append(&tmpString, plist[i]);
	}
	ov_string_freelist(plist);
	UA_NodeId tmpNodeId = UA_NODEID_NUMERIC(0, UA_NS0ID_HASTYPEDEFINITION);
	UA_NodeId tmpNodeId2 = UA_NODEID_NUMERIC(0, UA_NS0ID_HASCOMPONENT);
	for (size_t i = 0; i < newNode->referencesSize; i++){
		if (UA_NodeId_equal(&newNode->references[i].referenceTypeId, &tmpNodeId)){
			newNode->references[i].targetId = UA_EXPANDEDNODEID_NUMERIC(0, UA_NS0ID_PROPERTYTYPE);
			continue;
		}
		OV_STRING tmpString2 = NULL;
		copyOPCUAStringToOV(newNode->references[i].targetId.nodeId.identifier.string, &tmpString2);
		if (ov_string_compare(tmpString, tmpString2) == OV_STRCMP_EQUAL &&
			UA_NodeId_equal(&newNode->references[i].referenceTypeId, &tmpNodeId2) &&
			newNode->references[i].isInverse == UA_TRUE){
			newNode->references[i].referenceTypeId = UA_NODEID_NUMERIC(0, UA_NS0ID_HASPROPERTY);
		}
		ov_string_setvalue(&tmpString2, NULL);
	}
	ov_string_setvalue(&tmpString, NULL);
	UA_NodeId_deleteMembers(&tmpNodeId);
	UA_NodeId_deleteMembers(&tmpNodeId2);

	*opcuaNode = newNode;
	return UA_STATUSCODE_GOOD;
}

