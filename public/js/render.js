
// `reset` resets the state of the page to the default: the text field
// is usable, and there is no result div.
function reset() {
    input = $(":input").attr("disabled", false);
    $("#result").empty()
}

function render(parentDiv, tree) {
    // TODO: pull into a `renderNode()` or something
    if (tree.defn !== undefined) {

        var link = $('<a />',{
            text: tree.word,
            href: "http://dict.leo.org/englisch-deutsch/" + tree.word
        });
        $(parentDiv).append(link);

        var lines = tree.defn.split("\n");
        if (lines.length > 1) {
            var defn = lines[lines.length - 1];
            $(parentDiv).append($('<pre>').text(defn));
        }
    }

    if (tree.prefix !== undefined && tree.suffix !== undefined) {
        var tbl = $(document.createElement('table'));

        var tr = $(document.createElement('tr'));
        $(tbl).append(tr);

        var leftChild = $(document.createElement('td'));
        $(tr).append(leftChild);

        var rightChild = $(document.createElement('td'));
        $(tr).append(rightChild);

        render(leftChild, tree.prefix);
        render(rightChild, tree.suffix);

        $(parentDiv).append(tbl);
    }
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
    $.post(path)
        .fail(function(xhr, status, err) {
            $("#result").text("Error: " + status)
        })
        .done(function( data ) {
            reset();
            var tree = $( document.createElement('div') ).attr("id", "tree");
            $("#result").append(tree);

            var unmarshaledData = JSON.parse(data);
            render(tree, unmarshaledData);
        })
        .always(function() {
            input.attr("disabled", false);
        });
});

$(document).ready(function() {
    reset();
});
