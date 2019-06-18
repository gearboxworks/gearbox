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

#ifndef OV_COMPILE_LIBRARY_lifeCycleEntryOPCUAInterface
#define OV_COMPILE_LIBRARY_lifeCycleEntryOPCUAInterface
#endif

#include "lifeCycleEntryOPCUAInterface.h"
#include "libov/ov_macros.h"
#include "ksbase.h"
#include "opcua.h"
#include "opcua_helpers.h"
#include "NoneTicketAuthenticator.h"
#include "libov/ov_path.h"
#include "libov/ov_memstack.h"
#include "ks_logfile.h"
#include "nodeset_lifeCycleEntry.h"
#include "ua_lifeCycleEntry_generated.h"
#include "ua_lifeCycleEntry_generated_handling.h"


extern OV_INSTPTR_lifeCycleEntryOPCUAInterface_interface pinterface;

OV_DLLFNCEXPORT UA_StatusCode lifeCycleEntryOPCUAInterface_interface_ovLifeCycleEntryNodeToOPCUA(
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
	((UA_VariableNode*)newNode)->valueRank = -1;	/*	scalar	*/


	// value
	OV_ELEMENT tmpPart;
	tmpPart.elemtype = OV_ET_NONE;
	OV_ELEMENT tmpParrent;
	tmpParrent.pobj = pobj;
	tmpParrent.elemtype = OV_ET_OBJECT;
	UA_LifeCycleEntry tmpLifeCycleEntry;
	UA_LifeCycleEntry_init(&tmpLifeCycleEntry);
	do {
		ov_element_getnextpart(&tmpParrent, &tmpPart, OV_ET_VARIABLE);
		if (tmpPart.elemtype == OV_ET_NONE)
			break;
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "CreatingInstanceIdString") == OV_STRCMP_EQUAL){
			if (*(OV_STRING*)tmpPart.pvalue != NULL)
				tmpLifeCycleEntry.creatingInstance.idSpec = UA_String_fromChars(*(OV_STRING*)tmpPart.pvalue);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "CreatingInstanceIdType") == OV_STRCMP_EQUAL){
			tmpLifeCycleEntry.creatingInstance.idType = *(UA_UInt32*)tmpPart.pvalue;
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "WritingInstanceIdString") == OV_STRCMP_EQUAL){
			if (*(OV_STRING*)tmpPart.pvalue != NULL)
				tmpLifeCycleEntry.writingInstance.idSpec = UA_String_fromChars(*(OV_STRING*)tmpPart.pvalue);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "WritingInstanceIdType") == OV_STRCMP_EQUAL){
			tmpLifeCycleEntry.writingInstance.idType = *(UA_UInt32*)tmpPart.pvalue;
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "Data") == OV_STRCMP_EQUAL){
			ov_AnyToVariant((OV_ANY*)tmpPart.pvalue, &tmpLifeCycleEntry.data.value);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "TimeStamp") == OV_STRCMP_EQUAL){
			tmpLifeCycleEntry.data.sourceTimestamp = ov_ovTimeTo1601nsTime(*(OV_TIME*)tmpPart.pvalue);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "Subject") == OV_STRCMP_EQUAL){
			if (*(OV_STRING*)tmpPart.pvalue != NULL)
				tmpLifeCycleEntry.subject = UA_String_fromChars(*(OV_STRING*)tmpPart.pvalue);
			continue;
		}
		if (ov_string_compare(tmpPart.elemunion.pvar->v_identifier, "EventClass") == OV_STRCMP_EQUAL){
			if (*(OV_STRING*)tmpPart.pvalue != NULL)
				tmpLifeCycleEntry.eventClass = UA_String_fromChars(*(OV_STRING*)tmpPart.pvalue);
			continue;
		}
		tmpLifeCycleEntry.id =  atoi(tmpPart.elemunion.pvar->v_identifier);
	} while(TRUE);


	((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->type = &UA_LIFECYCLEENTRY[UA_LIFECYCLEENTRY_LIFECYCLEENTRY];
	((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data = UA_LifeCycleEntry_new();
	if (!((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data){
		result = UA_STATUSCODE_BADOUTOFMEMORY;
		return result;
	}
	((UA_VariableNode*)newNode)->value.data.value.hasValue = TRUE;
	((UA_VariableNode*)newNode)->valueSource = UA_VALUESOURCE_DATA;
	UA_LifeCycleEntry_copy(&tmpLifeCycleEntry, ((UA_Variant*)&((UA_VariableNode*)newNode)->value.data.value.value)->data);
	UA_LifeCycleEntry_deleteMembers(&tmpLifeCycleEntry);


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
	// dataType
	((UA_VariableNode*)newNode)->dataType = UA_NODEID_NUMERIC(pinterface->v_modelnamespace.index, UA_NSLIFECYCLEENTRYID_LIFECYCLEENTRY);


	// References
	addReference(newNode);
	UA_NodeId tmpNodeId = UA_NODEID_NUMERIC(0, UA_NS0ID_HASTYPEDEFINITION);
	for (size_t i = 0; i < newNode->referencesSize; i++){
		if (UA_NodeId_equal(&newNode->references[i].referenceTypeId, &tmpNodeId)){
			newNode->references[i].targetId = UA_EXPANDEDNODEID_NUMERIC(0, UA_NS0ID_PROPERTYTYPE);
		}
		if (newNode->references[i].isInverse == TRUE){
			ov_memstack_lock();
			result = opcua_nodeStoreFunctions_resolveNodeIdToPath(newNode->references[i].targetId.nodeId, &path);
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

			if (Ov_CanCastTo(lifeCycleEntry_LifeCycleArchive, pobj)){
				OV_STRING tmpOVString = NULL;
				copyOPCUAStringToOV(newNode->references[i].targetId.nodeId.identifier.string, &tmpOVString);
				ov_string_append(&tmpOVString, "|||LifeCyleEntries");
				UA_String tmpUAString = UA_String_fromChars(tmpOVString);
				UA_String_copy(&tmpUAString, &(newNode->references[i].targetId.nodeId.identifier.string));
			}
		}
	}
	UA_NodeId_deleteMembers(&tmpNodeId);


	*opcuaNode = newNode;
	return UA_STATUSCODE_GOOD;
}

