const Formatter = {

    date: function (s, format) {
        if(format)
            console.log("TODO: Formatter.date() with format: " + format);
        return new Date(Date.parse(s)).toLocaleDateString();
    }

}

