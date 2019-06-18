
/******************************************************************************
*
*   FILE
*   ----
*   MailBox.c
*
*   History
*   -------
*   2018-03-06   File created
*
*******************************************************************************
*
*   This file is generated by the 'acplt_builder' command
*
******************************************************************************/


#ifndef OV_COMPILE_LIBRARY_kbuslib
#define OV_COMPILE_LIBRARY_kbuslib
#endif


#include "kbuslib.h"
#include "libov/ov_macros.h"
#include "libov/ov_path.h"
#include "kbusl.h"

#define timespanToDouble(span)	\
	((OV_DOUBLE)(span)->secs) + ((OV_DOUBLE)(span)->usecs)/1000000.0

OV_DLLFNCEXPORT OV_RESULT kbuslib_MailBox_ByteAddressOut_set(
    OV_INSTPTR_kbuslib_MailBox          pobj,
    const OV_UINT  value
) {
    pobj->v_ByteAddressOut = value;
    return OV_ERR_OK;
}


OV_DLLFNCEXPORT OV_RESULT kbuslib_MailBox_busy_set(
    OV_INSTPTR_kbuslib_MailBox          pobj,
    const OV_BOOL  value
) {

	if(value){
		return kbuslib_MailBox_occupy(pobj, Ov_PtrUpCast(ov_object, pobj));
	} else {
		return kbuslib_MailBox_free(pobj, Ov_PtrUpCast(ov_object, pobj));
	}

}

/*
 * call ov_memstack_lock()/unlock() outside of this function
 */
OV_DLLFNCEXPORT OV_STRING* kbuslib_MailBox_queue_get(
    OV_INSTPTR_kbuslib_MailBox          pobj,
    OV_UINT *pveclen
) {
	*pveclen = pobj->v_occupier.len;
	OV_STRING* tmpStr = ov_memstack_alloc(*pveclen*sizeof(OV_STRING));
	if(!tmpStr){
		*pveclen = 0;
		return (OV_STRING*)0;
	}

	for(int i = 0; i<*pveclen; i++){
		tmpStr[i] = ov_path_getcanonicalpath(pobj->v_occupier.ppObj[i], 2);
	}

    return tmpStr;
}

OV_DLLFNCEXPORT void kbuslib_MailBox_startup(
	OV_INSTPTR_ov_object 	pobj
) {
    /*    
    *   local variables
    */
    OV_INSTPTR_kbuslib_MailBox pinst = Ov_StaticPtrCast(kbuslib_MailBox, pobj);

    /* do what the base class does first */
    kbuslib_Clamp_startup(pobj);

    /* do what */
    pinst->v_busy = FALSE;
	pinst->v_state = MB_idle;
	pinst->v_readBuffer = (OV_BYTE_VEC){.veclen=0, .value=NULL};
	pinst->v_writeBuffer = (OV_BYTE_VEC){.veclen=0, .value=NULL};
	pinst->v_occupier.ppObj = NULL;
	pinst->v_occupier.len = 0;

    return;
}

OV_DLLFNCEXPORT void kbuslib_MailBox_shutdown(
	OV_INSTPTR_ov_object 	pobj
) {
    /*
    *   local variables
    */
    OV_INSTPTR_kbuslib_MailBox pinst = Ov_StaticPtrCast(kbuslib_MailBox, pobj);

    /* do what */
	Ov_SetDynamicVectorLength(&pinst->v_readBuffer, 0 , BYTE);
	Ov_SetDynamicVectorLength(&pinst->v_writeBuffer, 0 , BYTE);
	Ov_HeapFree(pinst->v_occupier.ppObj);
	pinst->v_occupier.ppObj = NULL;
	pinst->v_occupier.len = 0;

    /* set the object's state to "shut down" */
    fb_functionblock_shutdown(pobj);

    return;
}

OV_DLLFNCEXPORT OV_RESULT kbuslib_MailBox_readwrite(
		OV_INSTPTR_kbuslib_MailBox	pobj,
		OV_BYTE* pMBoxPAE,
		OV_BYTE* pMBoxPAA,
		OV_BOOL* writeBack
	) {
	OV_UINT crc = 0x0;
	OV_UINT tmpCrc = 0x0;


	//if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80))
		//ov_logfile_debug("%s: in:  %02x %02x : %02x %02x %02x %02x", pobj->v_identifier, pMBoxPAE[0], pMBoxPAE[1],pMBoxPAE[2],pMBoxPAE[3],pMBoxPAE[4],pMBoxPAE[5]);


	switch((enum kbuslib_MBState)pobj->v_state){
	case MB_idle:
	case MB_finished:
	case MB_errorState:
		return OV_ERR_OK;

	case MB_startSend:
		//ov_logfile_debug("%s: start sending", pobj->v_identifier);
		pobj->v_writePos = 0;
		pobj->v_StatusByte = pMBoxPAE[1];

		pMBoxPAA[0] = 0x41; // Opcode 41=DLD_START
		pMBoxPAA[1] ^= 0x80; // toggle control bit
		pMBoxPAA[2] = 0x0;
		pMBoxPAA[3] = 0xFE; // transfer type byte array
		pMBoxPAA[4] = pobj->v_writeBuffer.veclen;
		pMBoxPAA[5] = 0x80; //TODO
		*writeBack = TRUE;
		pobj->v_state = MB_sending;
		break;

	case MB_sending:
		//ov_logfile_debug("%s: sending", pobj->v_identifier);
		if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80)){ // check for toggle bit
			pobj->v_StatusByte = pMBoxPAE[1];
			if(pMBoxPAE[0]!=pMBoxPAA[0] || (pMBoxPAE[1]&0x7f)){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}

			if((pobj->v_writePos<pobj->v_writeBuffer.veclen)){
				pMBoxPAA[0] = 0x42; // Opcode 42=DLD_CONT
				pMBoxPAA[1] ^= 0x80;

				for(OV_UINT i = 0; i<4; i++){
					if((pobj->v_writePos+i)<pobj->v_writeBuffer.veclen){
						pMBoxPAA[2+i] = pobj->v_writeBuffer.value[pobj->v_writePos+i];
					} else {
						pMBoxPAA[2+i] = 0x0;
					}
				}
				pobj->v_writePos += 4;
			} else {
				pMBoxPAA[0] = 0x43; // Opcode 43=DLD_END
				pMBoxPAA[1] ^= 0x80;
				for(int i = 0; i<pobj->v_writeBuffer.veclen;i++){
					crc += pobj->v_writeBuffer.value[i];
				}
				for(int i = 0; i<4; i++){
					pMBoxPAA[2+i] = crc >> (8*i);
				}
				pobj->v_state = MB_startRead;
			}
			*writeBack = TRUE;
		}
		break;

	case MB_stopSend:
		break;

	case MB_startRead:
		//ov_logfile_debug("%s: start reading", pobj->v_identifier);
		if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80)){ // check for toggle bit
			pobj->v_StatusByte = pMBoxPAE[1];
			if(pMBoxPAE[0]!=0x43 || (pMBoxPAE[1]&0x7f)){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}
			pobj->v_readPos = 0;

			pMBoxPAA[0] = 0x41;
			pMBoxPAA[1] ^= 0x80;
			pMBoxPAA[2] = 0x0;
			pMBoxPAA[3] = 0xfe;
			pMBoxPAA[4] = 0x0;
			pMBoxPAA[5] = 0xc0;
			*writeBack = TRUE;
			pobj->v_state = MB_startRead2;
		}
		break;

	case MB_startRead2:
		//ov_logfile_debug("%s: start reading 2", pobj->v_identifier);
		if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80)){ // check for toggle bit
			pobj->v_StatusByte = pMBoxPAE[1];
			if(pMBoxPAE[0]!=pMBoxPAA[0] || (pMBoxPAE[1]&0x7f)){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}
			Ov_SetDynamicVectorLength(&pobj->v_readBuffer, pMBoxPAE[4], BYTE);
			pMBoxPAA[0] = 0x42;
			pMBoxPAA[1] ^= 0x80;
			pMBoxPAA[2] = 0x0;
			pMBoxPAA[3] = 0x0;
			pMBoxPAA[4] = 0x0;
			pMBoxPAA[5] = 0x0;
			*writeBack = TRUE;
			pobj->v_state = MB_reading;
		}
		break;

	case MB_reading:
		//ov_logfile_debug("%s: reading", pobj->v_identifier);
		if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80)){ // check for toggle bit
			pobj->v_StatusByte = pMBoxPAE[1];
			if(pMBoxPAE[0]!=pMBoxPAA[0] || (pMBoxPAE[1]&0x7f)){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}

			for(int i= 0; i<4 && pobj->v_readPos+i < pobj->v_readBuffer.veclen; i++){
				pobj->v_readBuffer.value[pobj->v_readPos+i] = pMBoxPAE[2+i];
			}
			pobj->v_readPos += 4;
			if(pobj->v_readPos<pobj->v_readBuffer.veclen){
				pMBoxPAA[0] = 0x42;
			} else {
				pMBoxPAA[0] = 0x43;
				pobj->v_state = MB_stopRead;
			}
			pMBoxPAA[1] ^= 0x80;
			pMBoxPAA[2] = 0x0;
			pMBoxPAA[3] = 0x0;
			pMBoxPAA[4] = 0x0;
			pMBoxPAA[5] = 0x0;
			*writeBack = TRUE;
		}
		break;
	case MB_stopRead:
		//ov_logfile_debug("%s: stop reading", pobj->v_identifier);
		if((pMBoxPAE[1]&0x80) != (pobj->v_StatusByte&0x80)){ // check for toggle bit
			pobj->v_StatusByte = pMBoxPAE[1];
			if(pMBoxPAE[0]!= 0x43 || (pMBoxPAE[1]&0x7f)){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}

			tmpCrc = pMBoxPAE[2] + (pMBoxPAE[3] << 8) + (pMBoxPAE[4] << 16) + (pMBoxPAE[5] << 24);
			for(int i = 0; i<pobj->v_readBuffer.veclen; i++){
				crc += pobj->v_readBuffer.value[i];
			}
			if(crc != tmpCrc){
				pobj->v_state = MB_errorState;
				return OV_ERR_GENERIC;
			}
			pobj->v_state = MB_finished;
		}
		break;
	}

	//if(*writeBack)
		//ov_logfile_debug("%s: out: %02x %02x %02x %02x %02x %02x", pobj->v_identifier, pMBoxPAA[0], pMBoxPAA[1],pMBoxPAA[2],pMBoxPAA[3],pMBoxPAA[4],pMBoxPAA[5]);

    return OV_ERR_OK;
}

static OV_RESULT ov_queueAddElement(OV_QUEUE* pQueue, const OV_INSTPTR pobj){

	for(int i=0; i<pQueue->len; i++){
		if(pQueue->ppObj[i]==pobj)
			return OV_ERR_OK; // already in queue
	}

	pQueue->ppObj = Ov_HeapRealloc(pQueue->ppObj, (pQueue->len+1)*sizeof(OV_INSTPTR));
	if(!pQueue->ppObj){
		pQueue->len = 0;
		ov_logfile_error("failed to reallocate memory for queue");
		return OV_ERR_DBOUTOFMEMORY;
	}
	pQueue->ppObj[pQueue->len++] = pobj;
	return OV_ERR_OK;
}

static OV_RESULT ov_queuePopElement(OV_QUEUE* pQueue){

	OV_INSTPTR* ppObj = NULL;

	if(pQueue->len<2){ // 0 or 1; free queue
		pQueue->len = 0;
		pQueue->ppObj = NULL;
		return OV_ERR_OK;
	}

	ppObj = Ov_HeapMalloc((pQueue->len-1)*sizeof(OV_INSTPTR));
	if(!ppObj){
		ov_logfile_error("failed to allocate memory for queue");
		return OV_ERR_DBOUTOFMEMORY;
	}

	memcpy(ppObj, pQueue->ppObj+1, (pQueue->len-1)*sizeof(OV_INSTPTR)); // copy old list from 2nd element
	free(pQueue->ppObj);
	pQueue->ppObj = ppObj;
	pQueue->len--;
	return OV_ERR_OK;
}

// TODO implement more sophisticated occupy scheme
OV_DLLFNCEXPORT OV_RESULT kbuslib_MailBox_occupy(
	OV_INSTPTR_kbuslib_MailBox	pMBox,
	const OV_INSTPTR_ov_object		pobj
	) {

	OV_QUEUE* queue = &pMBox->v_occupier;
	OV_RESULT result;
	OV_TIME now;
	OV_TIME_SPAN diff;

	if(!queue->len){
		if(Ov_Fail(result=ov_queueAddElement(queue, pobj))){
			return result;
		}
		pMBox->v_busy = TRUE;
		return OV_ERR_OK;
	}
	if(queue->ppObj[0]==pobj){ // check if pobj is at the top
		pMBox->v_busy = TRUE;
		return OV_ERR_OK;
	}

	// check for timeout
	if(!pMBox->v_busy){
		ov_time_gettime(&now);
		ov_time_diff(&diff, &now, &pMBox->v_lastFree);
		if(timespanToDouble(&pMBox->v_queueTimeout)<timespanToDouble(&diff)){
			if(Ov_Fail(result=ov_queuePopElement(queue))){
				return result;
			}
			ov_time_gettime(&pMBox->v_lastFree);
		}
	}

	if(queue->ppObj[0]!=pobj){ // check if pobj is at the top now
		ov_queueAddElement(queue, pobj);
		return OV_ERR_NOACCESS;
	}

	pMBox->v_busy = TRUE;
	return OV_ERR_OK;
}

OV_DLLFNCEXPORT OV_RESULT kbuslib_MailBox_free(
	OV_INSTPTR_kbuslib_MailBox	pMBox,
	const OV_INSTPTR_ov_object		pobj
	) {

	OV_QUEUE* queue = &pMBox->v_occupier;
	OV_RESULT result;

	if(!queue->len) { // nothing to free
		pMBox->v_busy = FALSE;
		return OV_ERR_OK;
	}

	// check if we occupy the mailbox
	if(queue->ppObj[0]!=pobj &&
			pobj!=Ov_PtrUpCast(ov_object, pMBox)){ // override from mailbox itself
		return OV_ERR_NOACCESS;
	}

	if(Ov_Fail(result=ov_queuePopElement(queue))){
		return result;
	}

	ov_time_gettime(&pMBox->v_lastFree);
	pMBox->v_busy = FALSE;

    return OV_ERR_OK;
}
