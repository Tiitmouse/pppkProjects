// This file is part of the `@domagojpa/oib-validation` package.
// Licensed under the MIT License. See https://github.com/domagojpa/oib-validation/blob/main/LICENSE for details.
export function isOibValid(input: string): boolean {
    const oib = input.toString();

    if (oib.match(/^\d{11}$/) === null) {
        return false;
    }

    let calculated = 10;

    for (const digit of oib.substring(0, 10)) {
        calculated += parseInt(digit);
        calculated %= 10;

        if (calculated === 0) {
            calculated = 10;
        }

        calculated *= 2;
        calculated %= 11;
    }

    let check = 11 - calculated;
    
    if (check === 10) {
        check = 0;
    }

    return check === parseInt(oib[10]);
}