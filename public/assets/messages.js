(function(d){
    var showXML = function(){
        this.textContent = this.dataset.xmlDate;
    };
    var showFrom = function(){
        this.textContent = this.dataset.fromDate;
    };
    var bindTimeElementHandlers = function(element){
        element.onmouseover = showXML;
        element.onmouseout  = showFrom;
        element.onmouseup   = showFrom;
    }

    var updateTimes = function(){
        var els = d.getElementsByClassName("time");
        for(var i=0; i<els.length; i++){
            els[i].dataset.xmlDate = els[i].textContent;
            els[i].dataset.fromDate = moment(els[i].innerText).fromNow();
            els[i].textContent = els[i].dataset.fromDate;
            bindTimeElementHandlers(els[i]);
        }
    };

    d.onreadystatechange = function(){
        if(this.readyState === "complete"){
            updateTimes();
        }
    }
})(document);
