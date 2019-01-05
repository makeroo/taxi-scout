import {clone} from "lodash";


export function cloneAndPatch(obj, path, newValue, separator='.') {
    let stack = Array.isArray(path) ? path : path.split(separator);
    let newObj = clone(obj);

    obj = newObj;

    while (stack.length > 1) {
        let property = stack.shift();
        let sub = clone(obj[property]);

        obj[property] = sub;
        obj = sub;
    }

    obj[stack.shift()] = newValue;

    return newObj;
}

/*
function resolve(path, obj=self, separator='.') {
    var properties = Array.isArray(path) ? path : path.split(separator)
    return properties.reduce((prev, curr) => prev && prev[curr], obj)
}

function updateObject(object, newValue, path){

    var stack = path.split('>');

    while(stack.length>1){
        object = object[stack.shift()];
    }

    object[stack.shift()] = newValue;

}
*/
