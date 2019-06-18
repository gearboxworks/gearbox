
/******************************************************************************
 *
 *   FILE
 *   ----
 *   Slave.c
 *
 *   History
 *   -------
 *   2015-04-29   File created
 *
 *******************************************************************************
 *
 *   This file is generated by the 'acplt_builder' command
 *
 ******************************************************************************/


#ifndef OV_COMPILE_LIBRARY_modbusTcpLib
#define OV_COMPILE_LIBRARY_modbusTcpLib
#endif


#include "modbusTcpLib.h"
#include "libov/ov_macros.h"
#include "ksbase_helper.h"
#include "libov/ov_memstack.h"
#include "libov/ov_result.h"

OV_DLLFNCEXPORT OV_RESULT modbusTcpLib_Slave_host_set(
		OV_INSTPTR_modbusTcpLib_Slave          pobj,
		const OV_STRING  value
) {
	OV_VTBLPTR_TCPbind_TCPChannel	pVtblChannel = NULL;
	if(ov_string_compare(value, pobj->v_host) != OV_STRCMP_EQUAL){
		if(pobj->p_channel.v_ConnectionState != KSBASE_CONNSTATE_CLOSED){
			Ov_GetVTablePtr(TCPbind_TCPChannel, pVtblChannel, &(pobj->p_channel));
			if(!pVtblChannel){
				ov_logfile_error("%s: could not get Vtable of channel", pobj->v_identifier);
			} else {
				pVtblChannel->m_CloseConnection(Ov_PtrUpCast(ksbase_Channel, &(pobj->p_channel)));
			}
		}
	}
	return ov_string_setvalue(&pobj->v_host,value);
}

OV_DLLFNCEXPORT OV_RESULT modbusTcpLib_Slave_port_set(
		OV_INSTPTR_modbusTcpLib_Slave          pobj,
		const OV_INT  value
) {
	OV_VTBLPTR_TCPbind_TCPChannel	pVtblChannel = NULL;
	if(value != pobj->v_port){
		if(pobj->p_channel.v_ConnectionState != KSBASE_CONNSTATE_CLOSED){
			Ov_GetVTablePtr(TCPbind_TCPChannel, pVtblChannel, &(pobj->p_channel));
			if(!pVtblChannel){
				ov_logfile_error("%s: could not get Vtable of channel", pobj->v_identifier);
			} else {
				pVtblChannel->m_CloseConnection(Ov_PtrUpCast(ksbase_Channel, &(pobj->p_channel)));
			}
		}
	}
	pobj->v_port = value;
	return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_RESULT modbusTcpLib_Slave_ack_set(
		OV_INSTPTR_modbusTcpLib_Slave          pobj,
		const OV_BOOL  value
) {
	OV_INSTPTR_modbusTcpLib_IOChannel	pIOChannel	=	NULL;
	if(value == TRUE){
		Ov_ForEachChildEx(ov_containment, pobj, pIOChannel, modbusTcpLib_IOChannel){
			if(Ov_GetParent(modbusTcpLib_errorChannels, pIOChannel) == pobj){
				pIOChannel->v_error = FALSE;
				pIOChannel->v_errorCode = 0;
				modbusTcpLib_IOChannel_setErrorText(pIOChannel);
				Ov_Unlink(modbusTcpLib_errorChannels, pobj, pIOChannel);
			}
		}
		pobj->v_error = FALSE;
		if(pobj->p_channel.v_ConnectionState & KSBASE_CONNSTATE_ERRORINDICATOR){
			pobj->p_channel.v_ConnectionState = KSBASE_CONNSTATE_CLOSED;
		}
	}
	pobj->v_ack = value;
	return OV_ERR_OK;
}

OV_DLLFNCEXPORT void modbusTcpLib_Slave_typemethod(
		OV_INSTPTR_fb_functionblock	pfb,
		OV_TIME						*pltc
) {
	/*
	 *   local variables
	 */
	OV_INSTPTR_modbusTcpLib_Slave pinst = Ov_StaticPtrCast(modbusTcpLib_Slave, pfb);
	OV_VTBLPTR_TCPbind_TCPChannel	pVtblChannel	=	NULL;
	OV_INSTPTR_modbusTcpLib_Request	pRequest		=	NULL;
	OV_INSTPTR_modbusTcpLib_Request	pNextRequest		=	NULL;
	OV_VTBLPTR_modbusTcpLib_Request	pVtblRequest		=	NULL;
	OV_STRING	tempString = NULL;
	OV_INSTPTR_modbusTcpLib_IOChannel	pPrevChannel	=	NULL;
	OV_INSTPTR_modbusTcpLib_AoRI	pAoRI	=	NULL;
	OV_INSTPTR_modbusTcpLib_ReadInputRegistersRequest	pReadInputRegistersRequest	=	NULL;
	OV_INSTPTR_modbusTcpLib_AoRO	pAoRO	=	NULL;
	OV_INSTPTR_modbusTcpLib_WriteMultipleRegistersRequest	pWriteRegistersRequest	=	NULL;
	OV_INSTPTR_modbusTcpLib_DI		pDI		=	NULL;
	OV_INSTPTR_modbusTcpLib_ReadDiscreteInputsRequest	pReadDIRequest	=	NULL;
	OV_INSTPTR_modbusTcpLib_DO		pDO		=	NULL;
	OV_INSTPTR_modbusTcpLib_WriteMultipleCoilsRequest		pWriteCoilsRequest	=	NULL;
	OV_INT		currUnitIdentifier	= 0xFF;
	OV_INT		nextUnitIdentifier	=	-1;
	OV_RESULT result;

	/*	check request timeouts	*/
	pRequest = Ov_DynamicPtrCast(modbusTcpLib_Request, Ov_GetFirstChild(ov_containment, &(pinst->p_requests)));
	while(pRequest){
		pNextRequest = Ov_DynamicPtrCast(modbusTcpLib_Request, Ov_GetNextChild(ov_containment, pRequest));
		if(ov_time_compare(&(pRequest->v_validTill), pltc) == OV_TIMECMP_BEFORE){
			Ov_DeleteObject(pRequest);
		}
		pRequest = pNextRequest;
	}

	switch(pinst->p_channel.v_ConnectionState){
	case KSBASE_CONNSTATE_CLOSED:
		Ov_GetVTablePtr(TCPbind_TCPChannel, pVtblChannel, &(pinst->p_channel));
		if(!pVtblChannel){
			ov_logfile_error("%s: could not get Vtable of channel", pinst->v_identifier);
		} else {
			ov_string_print(&tempString, "%i", pinst->v_port);
			pVtblChannel->m_OpenConnection(Ov_PtrUpCast(ksbase_Channel, &(pinst->p_channel)), pinst->v_host, tempString);
			ov_string_setvalue(&tempString, NULL);
		}
		break;
	case KSBASE_CONNSTATE_OPEN:
		/*	do the bus-request stuff --> do an interation for every type of IOChannel	*/
		/*	AoRIs	*/
		do{
			nextUnitIdentifier = -1;
			Ov_ForEachChildEx(ov_containment, pinst, pAoRI, modbusTcpLib_AoRI){
				if(Ov_GetParent(modbusTcpLib_requestToChannel, pAoRI) == NULL && Ov_GetParent(modbusTcpLib_toNextChannel, pAoRI) == NULL){
					/*	no Request running for current IOChannel	*/
					if(pAoRI->v_unitIdentifier == currUnitIdentifier && pAoRI->v_actimode == 1){
						if(!pReadInputRegistersRequest){
							result = Ov_CreateIDedObject(modbusTcpLib_ReadInputRegistersRequest, pReadInputRegistersRequest, &(pinst->p_requests), "ReadInputRegisters");
							if(Ov_Fail(result)){
								ov_logfile_error("%s: creation of Request failed with error: %s", pinst->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pReadInputRegistersRequest->v_unitIdentifier = pAoRI->v_unitIdentifier;
							result = Ov_Link(modbusTcpLib_requestToChannel, pReadInputRegistersRequest, pAoRI);
							if(Ov_Fail(result)){
								ov_logfile_error("%s: Linking of first Channel (%s) to Request (%s) failed with error: %s", pinst->v_identifier,
										pAoRI->v_identifier, pReadInputRegistersRequest->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRI);
						}
						Ov_GetVTablePtr(modbusTcpLib_Request, pVtblRequest, pReadInputRegistersRequest);
						if(!pVtblRequest){
							ov_logfile_error("%s: could not get Vtable of Request %s", pinst->v_identifier, pReadInputRegistersRequest->v_identifier);
							break;
						}
						if(pReadInputRegistersRequest->v_requestStartAddr == -1
								|| (pAoRI->v_address < pReadInputRegistersRequest->v_requestStartAddr
										&& pAoRI->v_address > ((pReadInputRegistersRequest->v_requestStartAddr + pReadInputRegistersRequest->v_requestedItems - 125)))
									|| ((pAoRI->v_address > pReadInputRegistersRequest->v_requestStartAddr)
										&& (pAoRI->v_address < (pReadInputRegistersRequest->v_requestStartAddr + 125)))){
							pVtblRequest->m_addItem(Ov_PtrUpCast(modbusTcpLib_Request, pReadInputRegistersRequest), pAoRI->v_address);
							if(pPrevChannel != Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRI)){
								result = Ov_Link(modbusTcpLib_toNextChannel, pPrevChannel, pAoRI);
								if(Ov_Fail(result)){
									ov_logfile_error("%s: Linking of Channel (%s) to previous Channel (%s) failed with error: %s", pinst->v_identifier,
											pAoRI->v_identifier, pPrevChannel->v_identifier, ov_result_getresulttext(result));
									break;
								}
								pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRI);
							}
						}
						if(pReadInputRegistersRequest->v_requestedItems >= 125){
							/*	125 registers per request tops
							 * do another run					*/
							if(nextUnitIdentifier == -1){
								nextUnitIdentifier = currUnitIdentifier;
							}
							pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pReadInputRegistersRequest));
							break;
						}
					} else {
						/*	this IOCHannel was not included in this request --> do another iteration with the right unitIdentifer	*/
						if(nextUnitIdentifier == -1 && pAoRI->v_actimode == 1){
							nextUnitIdentifier = pAoRI->v_unitIdentifier;
						}
					}
				}
			}
			if(nextUnitIdentifier != -1){
				currUnitIdentifier = nextUnitIdentifier;
			}
			if(pReadInputRegistersRequest && pVtblRequest){
				pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pReadInputRegistersRequest));
			}
		} while(nextUnitIdentifier != -1);

		/*	AoROs	*/
		do{
			nextUnitIdentifier = -1;
			Ov_ForEachChildEx(ov_containment, pinst, pAoRO, modbusTcpLib_AoRO){
				if(Ov_GetParent(modbusTcpLib_requestToChannel, pAoRO) == NULL && Ov_GetParent(modbusTcpLib_toNextChannel, pAoRO) == NULL){
					/*	no Request running for current IOChannel	*/
					if(pAoRO->v_unitIdentifier == currUnitIdentifier && pAoRO->v_actimode == 1){
						if(!pWriteRegistersRequest){
							result = Ov_CreateIDedObject(modbusTcpLib_WriteMultipleRegistersRequest, pWriteRegistersRequest, &(pinst->p_requests), "WriteRegisters");
							if(Ov_Fail(result)){
								ov_logfile_error("%s: creation of Request failed with error: %s", pinst->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pWriteRegistersRequest->v_unitIdentifier = pAoRO->v_unitIdentifier;
							result = Ov_Link(modbusTcpLib_requestToChannel, pWriteRegistersRequest, pAoRO);
							if(Ov_Fail(result)){
								ov_logfile_error("%s: Linking of first Channel (%s) to Request (%s) failed with error: %s", pinst->v_identifier,
										pAoRO->v_identifier, pWriteRegistersRequest->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRO);
						}
						Ov_GetVTablePtr(modbusTcpLib_Request, pVtblRequest, pWriteRegistersRequest);
						if(!pVtblRequest){
							ov_logfile_error("%s: could not get Vtable of Request %s", pinst->v_identifier, pWriteRegistersRequest->v_identifier);
							break;
						}
						if(pWriteRegistersRequest->v_requestStartAddr == -1
								|| (pAoRO->v_address < pWriteRegistersRequest->v_requestStartAddr
										&& pAoRO->v_address > ((pWriteRegistersRequest->v_requestStartAddr + pWriteRegistersRequest->v_requestedItems - 125)))
									|| ((pAoRO->v_address > pWriteRegistersRequest->v_requestStartAddr)
										&& (pAoRO->v_address < (pWriteRegistersRequest->v_requestStartAddr + 125)))){
							pVtblRequest->m_addItem(Ov_PtrUpCast(modbusTcpLib_Request, pWriteRegistersRequest), pAoRO->v_address);
							if(pPrevChannel != Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRO)){
								result = Ov_Link(modbusTcpLib_toNextChannel, pPrevChannel, pAoRO);
								if(Ov_Fail(result)){
									ov_logfile_error("%s: Linking of Channel (%s) to previous Channel (%s) failed with error: %s", pinst->v_identifier,
											pAoRO->v_identifier, pPrevChannel->v_identifier, ov_result_getresulttext(result));
									break;
								}
								pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pAoRO);
							}
						}
						if(pWriteRegistersRequest->v_requestedItems >= 125){
							/*	125 registers per request tops
							 * do another run					*/
							if(nextUnitIdentifier == -1){
								nextUnitIdentifier = currUnitIdentifier;
							}
							pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pWriteRegistersRequest));
							break;
						}
					} else {
						/*	this IOCHannel was not included in this request --> do another iteration with the right unitIdentifer	*/
						if(nextUnitIdentifier == -1 && pAoRO->v_actimode == 1){
							nextUnitIdentifier = pAoRO->v_unitIdentifier;
						}
					}
				}
			}
			if(nextUnitIdentifier != -1){
				currUnitIdentifier = nextUnitIdentifier;
			}
			if(pWriteRegistersRequest && pVtblRequest){
				pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pWriteRegistersRequest));
			}
		} while(nextUnitIdentifier != -1);

		/*	DIs	*/
		do{
			nextUnitIdentifier = -1;
			Ov_ForEachChildEx(ov_containment, pinst, pDI, modbusTcpLib_DI){
				if(Ov_GetParent(modbusTcpLib_requestToChannel, pDI) == NULL && Ov_GetParent(modbusTcpLib_toNextChannel, pDI) == NULL){
					/*	no Request running for current IOChannel	*/
					if(pDI->v_unitIdentifier == currUnitIdentifier && pDI->v_actimode == 1){
						if(!pReadDIRequest){
							result = Ov_CreateIDedObject(modbusTcpLib_ReadDiscreteInputsRequest, pReadDIRequest, &(pinst->p_requests), "ReadDiscreteInputs");
							if(Ov_Fail(result)){
								ov_logfile_error("%s: creation of Request failed with error: %s", pinst->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pReadDIRequest->v_unitIdentifier = pDI->v_unitIdentifier;
							result = Ov_Link(modbusTcpLib_requestToChannel, pReadDIRequest, pDI);
							if(Ov_Fail(result)){
								ov_logfile_error("%s: Linking of first Channel (%s) to Request (%s) failed with error: %s", pinst->v_identifier,
										pDI->v_identifier, pReadDIRequest->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pDI);
						}
						Ov_GetVTablePtr(modbusTcpLib_Request, pVtblRequest, pReadDIRequest);
						if(!pVtblRequest){
							ov_logfile_error("%s: could not get Vtable of Request %s", pinst->v_identifier, pReadDIRequest->v_identifier);
							break;
						}
						if(pReadDIRequest->v_requestStartAddr == -1
								|| (pDI->v_address < pReadDIRequest->v_requestStartAddr
										&& pDI->v_address > ((pReadDIRequest->v_requestStartAddr + pReadDIRequest->v_requestedItems - 2000)))
								|| ((pDI->v_address > pReadDIRequest->v_requestStartAddr)
										&& (pDI->v_address < (pReadDIRequest->v_requestStartAddr + 2000)))){
							pVtblRequest->m_addItem(Ov_PtrUpCast(modbusTcpLib_Request, pReadDIRequest), pDI->v_address);
							if(pPrevChannel != Ov_PtrUpCast(modbusTcpLib_IOChannel, pDI)){
								result = Ov_Link(modbusTcpLib_toNextChannel, pPrevChannel, pDI);
								if(Ov_Fail(result)){
									ov_logfile_error("%s: Linking of Channel (%s) to previous Channel (%s) failed with error: %s", pinst->v_identifier,
											pDI->v_identifier, pPrevChannel->v_identifier, ov_result_getresulttext(result));
									break;
								}
								pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pDI);
							}
						}
						if(pReadDIRequest->v_requestedItems >= 2000){
							/*	2000 bits per request tops
							 * do another run					*/
							if(nextUnitIdentifier == -1){
								nextUnitIdentifier = currUnitIdentifier;
							}
							pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pReadDIRequest));
							break;
						}
					} else {
						/*	this IOCHannel was not included in this request --> do another iteration with the right unitIdentifer	*/
						if(nextUnitIdentifier == -1 && pDI->v_actimode == 1){
							nextUnitIdentifier = pDI->v_unitIdentifier;
						}
					}
				}
			}
			if(nextUnitIdentifier != -1){
				currUnitIdentifier = nextUnitIdentifier;
			}
			if(pReadDIRequest && pVtblRequest){
				pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pReadDIRequest));
			}
		} while(nextUnitIdentifier != -1);

		/*	DOs	*/
		do{
			nextUnitIdentifier = -1;
			Ov_ForEachChildEx(ov_containment, pinst, pDO, modbusTcpLib_DO){
				if(Ov_GetParent(modbusTcpLib_requestToChannel, pDO) == NULL && Ov_GetParent(modbusTcpLib_toNextChannel, pDO) == NULL){
					/*	no Request running for current IOChannel	*/
					if(pDO->v_unitIdentifier == currUnitIdentifier && pDO->v_actimode == 1){
						if(!pWriteCoilsRequest){
							result = Ov_CreateIDedObject(modbusTcpLib_WriteMultipleCoilsRequest, pWriteCoilsRequest, &(pinst->p_requests), "WriteCoils");
							if(Ov_Fail(result)){
								ov_logfile_error("%s: creation of Request failed with error: %s", pinst->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pWriteCoilsRequest->v_unitIdentifier = pDO->v_unitIdentifier;
							result = Ov_Link(modbusTcpLib_requestToChannel, pWriteCoilsRequest, pDO);
							if(Ov_Fail(result)){
								ov_logfile_error("%s: Linking of first Channel (%s) to Request (%s) failed with error: %s", pinst->v_identifier,
										pDO->v_identifier, pWriteCoilsRequest->v_identifier, ov_result_getresulttext(result));
								break;
							}
							pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pDO);
						}
						Ov_GetVTablePtr(modbusTcpLib_Request, pVtblRequest, pWriteCoilsRequest);
						if(!pVtblRequest){
							ov_logfile_error("%s: could not get Vtable of Request %s", pinst->v_identifier, pWriteCoilsRequest->v_identifier);
							break;
						}
						if(pWriteCoilsRequest->v_requestStartAddr == -1
								|| (pDO->v_address < pWriteCoilsRequest->v_requestStartAddr
										&& pDO->v_address > ((pWriteCoilsRequest->v_requestStartAddr + pWriteCoilsRequest->v_requestedItems - 2000)))
								|| ((pDO->v_address > pWriteCoilsRequest->v_requestStartAddr)
										&& (pDO->v_address < (pWriteCoilsRequest->v_requestStartAddr + 2000)))){
							pVtblRequest->m_addItem(Ov_PtrUpCast(modbusTcpLib_Request, pWriteCoilsRequest), pDO->v_address);
							if(pPrevChannel != Ov_PtrUpCast(modbusTcpLib_IOChannel, pDO)){
								result = Ov_Link(modbusTcpLib_toNextChannel, pPrevChannel, pDO);
								if(Ov_Fail(result)){
									ov_logfile_error("%s: Linking of Channel (%s) to previous Channel (%s) failed with error: %s", pinst->v_identifier,
											pDO->v_identifier, pPrevChannel->v_identifier, ov_result_getresulttext(result));
									break;
								}
								pPrevChannel = Ov_PtrUpCast(modbusTcpLib_IOChannel, pDO);
							}
						}
						if(pWriteCoilsRequest->v_requestedItems >= 2000){
							/*	2000 bits per request tops
							 * do another run					*/
							if(nextUnitIdentifier == -1){
								nextUnitIdentifier = currUnitIdentifier;
							}
							pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pWriteCoilsRequest));
							break;
						}
					} else {
						/*	this IOCHannel was not included in this request --> do another iteration with the right unitIdentifer	*/
						if(nextUnitIdentifier == -1 && pDO->v_actimode == 1){
							nextUnitIdentifier = pDO->v_unitIdentifier;
						}
					}
				}
			}
			if(nextUnitIdentifier != -1){
				currUnitIdentifier = nextUnitIdentifier;
			}
			if(pWriteCoilsRequest && pVtblRequest){
				pVtblRequest->m_sendRequest(Ov_PtrUpCast(modbusTcpLib_Request, pWriteCoilsRequest));
			}
		} while(nextUnitIdentifier != -1);

		break;
	default:
		if(pinst->p_channel.v_ConnectionState & KSBASE_CONNSTATE_ERRORINDICATOR){
			pinst->v_error = TRUE;
		}
		break;
	}

	return;
}

OV_DLLFNCEXPORT OV_RESULT modbusTcpLib_Slave_constructor(
		OV_INSTPTR_ov_object 	pobj
) {
	/*
	 *   local variables
	 */
	OV_INSTPTR_modbusTcpLib_Slave	pinst = Ov_StaticPtrCast(modbusTcpLib_Slave, pobj);
	OV_INSTPTR_ov_domain			pParentDomain	=	NULL;
	OV_INSTPTR_modbusTcpLib_ModbusTcpManager	pManager	=	NULL;
	OV_RESULT    result;

	/* do what the base class does first */
	result = fb_functionblock_constructor(pobj);
	if(Ov_Fail(result))
		return result;

	/* do what */
	pParentDomain = Ov_GetParent(ov_containment, pobj);
	if(!pParentDomain){
		pParentDomain = Ov_StaticPtrCast(ov_domain, pobj->v_pouterobject);
	}
	if(pParentDomain){
		pManager = Ov_DynamicPtrCast(modbusTcpLib_ModbusTcpManager, pParentDomain);
		if(!pManager){
			return OV_ERR_BADPATH;
		}
	} else {
		return OV_ERR_GENERIC;
	}
	result = Ov_LinkPlaced(fb_tasklist, pManager, pinst, OV_PMH_END);
	if(Ov_Fail(result)){
		return result;
	}
	pinst->p_channel.v_CloseAfterSend = FALSE;
	pinst->p_channel.v_ConnectionTimeOut = pinst->v_timeout;
	pinst->p_channel.v_UnusedDataTimeOut = pinst->v_timeout;
	return Ov_Link(ksbase_AssocChannelDataHandler, &(pinst->p_channel), &(pinst->p_dispatcher));
}

OV_DLLFNCEXPORT OV_RESULT modbusTcpLib_Slave_rename(
		OV_INSTPTR_ov_object	pobj,
		OV_STRING				newid,
		OV_INSTPTR_ov_domain	pnewparent
) {
	OV_INSTPTR_modbusTcpLib_ModbusTcpManager	pManager	=	NULL;
	OV_INSTPTR_modbusTcpLib_Slave pinst = Ov_StaticPtrCast(modbusTcpLib_Slave, pobj);
	OV_INSTPTR_fb_task	pOldTask	=	NULL;
	OV_RESULT	result;
	pManager = Ov_DynamicPtrCast(modbusTcpLib_ModbusTcpManager, pnewparent);
	if(!pManager){
		return OV_ERR_BADPATH;
	}
	pOldTask = Ov_GetParent(fb_tasklist, pinst);
	Ov_Unlink(fb_tasklist, pOldTask, pinst);
	result = Ov_LinkPlaced(fb_tasklist, pManager, pinst, OV_PMH_END);
	if(Ov_Fail(result)){
		Ov_LinkPlaced(fb_tasklist, pOldTask, pinst, OV_PMH_END);
		return result;
	}
	return OV_ERR_OK;
}

