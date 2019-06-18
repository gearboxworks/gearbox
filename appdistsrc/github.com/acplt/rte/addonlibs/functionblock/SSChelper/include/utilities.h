/*
 * utilities.h
 *
 *  Created on: Aug 30, 2017
 *      Author: julian
 */

#ifndef UTILITIES_H_
#define UTILITIES_H_

OV_RESULT SSChelper_getObjectAndVarnameFromSetVariable(
		const OV_INSTPTR_ov_object pinst,
		const OV_STRING nameToCheck,
		OV_INSTPTR_ov_object *pTargetObj,
		OV_STRING *ptargetVarname
);

OV_RESULT SSChelper_setNamedVariable(const OV_INSTPTR_ov_object pTargetObj, const OV_STRING targetVarname, OV_ANY *value);

OV_RESULT SSChelper_getNamedVariable(const OV_INSTPTR_ov_object pTargetObj, const OV_STRING targetVarname, OV_ANY *value);

#endif /* UTILITIES_H_ */
