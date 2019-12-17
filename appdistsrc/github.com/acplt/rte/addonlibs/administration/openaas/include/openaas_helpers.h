/*
 * openaas_helpers.h
 *
 *  Created on: 07.02.2017
 *      Author: ubuntu
 */

#ifndef OPENAAS_HELPERS_H_
#define OPENAAS_HELPERS_H_

#include "openaas.h"
#include "jsonparsing.h"
#include "lifeCycleEntry_helpers.h"
#include "propertyValueStatement_helpers.h"

OV_RESULT decodeMSG(const SRV_String* str, SRV_msgHeader** header, void** srvStruct, SRV_service_t* srvType, SRV_encoding_t *encoding);
OV_RESULT encodeMSG(SRV_String** str, const SRV_msgHeader *header, const void* srvStruct, SRV_service_t srvType, SRV_encoding_t encoding);
OV_RESULT serviceValueToOVDataValue(OV_ANY* value, const SRV_extAny_t* serviceValue);
OV_RESULT OVDataValueToserviceValue(OV_ANY value, SRV_extAny_t* serviceValue);
AASStatusCode checkForEmbeddingAAS(IdentificationType aasId, IdentificationType objectId);
OV_BOOL getAASIdbyObjectPointer(OV_INSTPTR_openaas_aas pAAS, IdentificationType* pAASId);
OV_STRING getStatementsInJSON(OV_INSTPTR_openaas_aas paas);
#endif /* OPENAAS_HELPERS_H_ */
