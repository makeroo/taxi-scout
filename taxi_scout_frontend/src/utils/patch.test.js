import {cloneAndPatch} from "./patch";


it('clone and patch simple property', () => {
    const obj = {
        a: 2,
        b: 3,
        c: [4, 5]
    };

    const newObj = cloneAndPatch(obj, 'a', 6);

    expect(obj).toEqual({
        a: 2,
        b: 3,
        c: [4, 5]
    });

    expect(newObj).toEqual({
        a: 6,
        b: 3,
        c: [4, 5]
    });

    expect(obj.c).toBe(newObj.c);
});


it('clone and patch nested property', () => {
    const obj = {
        a: 2,
        b: 3,
        c: { x:4, y:5 }
    };

    const newObj = cloneAndPatch(obj, 'c.x', 6);

    expect(obj).toEqual({
        a: 2,
        b: 3,
        c: { x:4, y:5 }
    });

    expect(newObj).toEqual({
        a: 2,
        b: 3,
        c: { x:6, y:5 }
    });
});


it('clone and patch array', () => {
    const obj = {
        a: 2,
        b: 3,
        c: [4, 5]
    };

    const newObj = cloneAndPatch(obj, ['c', '1'], 6);

    expect(obj).toEqual({
        a: 2,
        b: 3,
        c: [4, 5]
    });

    expect(newObj).toEqual({
        a: 2,
        b: 3,
        c: [4, 6]
    });
});
