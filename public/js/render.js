
// `reset` resets the state of the page to the default: the text field
// is usable, and there is no result div.
function reset() {
    input = $(":input").attr("disabled", false);
    $("#result").empty()
}

function render(tree) {
    if (tree === undefined) {
	return undefined;
    }

    console.log(tree);
    var node = $('<div/>', {class: 'node'});

    if (tree.word !== undefined && tree.defns !== undefined) {
        var link = $('<a />',{
            text: tree.word,
            href: "http://dict.leo.org/englisch-deutsch/" + tree.word
        });
        $(node).append(link);

        var defn = $('<div/>', {class: 'defn'});
        $(node).append(defn);

        var n = (tree.defns.length < 3 ? tree.defns.length : 3);
        for (i = 0; i < n; i++) {
            var d = tree.defns[i];
            $(defn).append($('<pre>').text(d));
        }
    }

    var ul = $(document.createElement('ul'));
    $(node).append(ul);
    
    var leftEl = render(tree.prefix);
    var rightEl = render(tree.suffix);

    if (leftEl !== undefined) {
        $(ul).append(leftEl);
    }
    if (rightEl !== undefined) {
        $(ul).append(rightEl);
    }
    
    return node;
}

$( "#lookup" ).submit(function( event ) {
    event.preventDefault();

    reset();
    console.log("!");

    var $form = $( this ),
        input = $form.find( "input[name='q']" );
    var path = "/split/" + input.val();

    $("#result").text("Loading").show();

    input.attr("disabled", true);
    $.get(path)
        .fail(function(xhr, status, err) {
            $("#result").text("Error: " + status)
        })
        .done(function( data ) {
            reset();
            var tree = $( document.createElement('div') ).attr("id", "tree");
            $("#result").append(tree);

            var ul = $(document.createElement('ul'));
            var li = $(document.createElement('li'));
            ul.append(li);
            tree.append(ul);
        
            var unmarshaledData = JSON.parse(data);
            li.append(render(unmarshaledData));
        })
        .always(function() {
            input.attr("disabled", false);
        });
});

$(document).ready(function() {
    reset();
});
