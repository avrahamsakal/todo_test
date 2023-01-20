// ===Event handlers===
//   (Handler funcs aren't defined anonymously because
//    event handlers need to be re-attached every time
//    a new row gets added)

//$('input[type=text]').change(textInput_Change)
$('.fields').change(field_Change)
const field_Change = (e) => {
    update(e)
}
$('.delete').click(deleteButton_Click)
const deleteButton_Click = (e) => {
    del(e)
}

// ===Support functions for base CRUD operations===

const getEntityFromInputName = (e) => {
    // Read key and value pair from form input
    const [fullKey, v] = [$(e).attr('name'), $(e).val()];
    const [entityName, ...args] = fullKey.replaceAll(']', '').split('[');
    const [entityId, entity] = [args[1] || '', args[2] ? {k:v} : {}];
    return [entityName, entityId, entity];
}

const getParentRow = (e) => {
    return $(e).filter('div:parents:first')
}


// ===CRUD operations===

const create = (e) => {
    return update(e);
}

const read = (e, successFunc) => {
    const [entityName, entityId, entity] = getEntityFromInputName(e);

    // Post new value to the place it says it's from
    $.post('/' + entityName + '/' + entityId, entity).success((data, textStatus, jqXHR) => {
        if (successFunc) successFunc(e, data, textStatus, jqXHR)
    }).fail((err) => {
        // @TODO error handling
    })
    // return entity

    /*const successFuncExample = (e, data, textStatus, jqXHR) => {
        // This is just an example success handler for a read operation
    }*/
}

const update = (e) => {
    const [entityName, entityId, entity] = getEntityFromInputName(e);

    // Post new value to the place it says it's from
    // @TODO: Technically this should be $.ajax({'type':'PATCH'})
    $.post('/' + entityName + '/' + entityId, entity).success((data, textStatus, jqXHR) => {
        // @TODO check HTTP status is 200 OK or updated status code
        $(e).attr('name', $(e).attr('name').replace(
            '[0]', '['+ (data.id || 0).toString() + ']'
        ));
        // Need last_inserted_id because this might be a create
    }).fail((err) => {
        // @TODO error handling
    })

    // If entityId == 0 (i.e. a create operation), need to
    //  clone the current row and append it to current row
    //  and hook up all the missing event handlers so that
    //  we can keep adding new rows automatically ad inf.
    const row = getParentRow(e);
    const newRow = $(row).after(row);
    newRow.filter('.field').change(field_Change)
    newRow.filter('.delete').click(deleteButton_Click);

    // return null because async
}

const del = (e) => {
    getParentRow(e).remove() // Delete entire row
    
    const [entityName, entityId, _] = getEntityFromInputName(e);
    // Execute DELETE verb against where the input name says to
    $.ajax({
        type: 'DELETE',
        url: '/' + entityName + '/' + entityId,
        success: (data, textStatus, jqXHR) => {
            // @TODO check HTTP status is 200 OK,
            //   undelete row on failure, very easy
            //   to do in jQuery using getParentRow(e)            
        }, error: (err) => { /* @TODO error handling */ }
    })
}
