function GetTextCompareClass(val1, val2) {
    if (val1 > val2) {
        return "text-end text-success"
    } else if (val1 === val2) {
        return "text-end"
    } else {
        return "text-end text-danger"
    }
};

function GetDeltaText(val1, val2, formatOptions, addDollarSign) {
    if (val1 > val2) {
        return `- ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val1 - val2))}`
    } else if (val1 === val2) {
        return "-"
    } else {
        return `+ ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val1 - val2))}`
    }
}

function GetDeltaTextForSingleValue(val, formatOptions, addDollarSign) {
    if (val > 0) {
        return `+ ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val))}`
    } else if (val < 0) {
        return `- ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val))}`
    } else {
        return "-"
    }
}

function GetDeltaMonthText(val1, val2) {
    if (val1 > val2) {
        return `- ${val1 - val2}`
    } else if (val1 === val2) {
        return "-"
    } else {
        return `+ ${val1 - val2}`
    }
}

export {GetDeltaMonthText, GetDeltaText, GetDeltaTextForSingleValue, GetTextCompareClass};