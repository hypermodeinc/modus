import { registerFunction } from "./host";
import { HMFunction } from "../protos/HMFunction";
import { HMParameter } from "../protos/HMParameter";
import { HMType } from "../protos/HMType";
import { Protobuf } from "as-proto/assembly";

export function register(resolver: string, name: string, paramNames: string[], paramTypes: HMType[], returnType: HMType) : void {

    if (paramNames.length != paramTypes.length) {
        throw new Error("paramNames and paramTypes must be the same length");
    }

    const params = new Array<HMParameter>(paramNames.length);
    for (let i = 0; i < paramNames.length; i++) {
        const name = paramNames[i];
        const type = paramTypes[i];
        if (type == HMType.VOID) {
            throw new Error("paramTypes cannot contain VOID");
        }

        params[i] = new HMParameter(name, type);
    }

    const func = new HMFunction(resolver, name, params, returnType);
    const msg = Protobuf.encode(func, HMFunction.encode);
    registerFunction(msg.buffer);
}
