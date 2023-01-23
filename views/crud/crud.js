// ===Event handlers===
//   (Handler funcs aren't defined anonymously because
//    event handlers need to be re-attached every time
//    a new row gets added)

/*const field_Change = (e) => { update(e); }
const deleteButton_Click = (e) => { return del(e); }
const submitButton_Click = (e) => { $(e).parents('form:first').submit(); return true; }
$(() => {
    $('.field').change(field_Change($(this)));
    $('.delete').click(deleteButton_Click($(this)));
    $('.submit').click(submitButton_Click($(this)));
});*/

// ===Support functions for base CRUD operations===

const getEntityFromInputName = (e) => {
    // Read key and value pair from form input
    const tagName = $(e).prop('tagName');
    console.log("Tag name:", tagName);
    const [fullKey, v] = (tagName === 'INPUT') ?
        [$(e).attr('name') || '', $(e).val() || ''] :
        [$(e).data('name') || '', ''];
    console.log("Key:", fullKey, "Value:", v);
    const [entityName, ...args] = fullKey.replaceAll(']', '').split('[');
    console.log("Name:", entityName, "Args:", args);
    const [entityId, k] = [args[0] || '0', (args.length > 0) ? args[1].toString() : ''];
    const entity = {};
    if (k !== '') entity[k] = v;
    console.log("ID:", entityId, "Args:", args, "Data:", entity);
    return [entityName, entityId, entity];
}

const getParentRow = (e) => {
    return $(e).parents('div:first');
}


// ===CRUD operations===

const create = (e, isCollection) => {
    return update(e, isCollection);
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
    const isCollection = ($(e).data('collection') !== undefined)

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

    if (isCollection && (entityId === 0 || entityId === '0')) {
        // (i.e. a create operation), need to
        //  clone the current row and append it to current row
        //  and hook up all the missing event handlers so that
        //  we can keep adding new rows automatically ad inf.
        const row = getParentRow(e);
        const newRow = $(row).clone();
        $(newRow).children('input').val('');
        $(row).after(newRow);
        $(newRow).fadeOut(100).fadeIn(100).fadeOut(100).fadeIn(100); // Flash new row

        //newRow.filter('.field').change(field_Change)
        //newRow.filter('.delete').click(deleteButton_Click);
    }
    // return null because async
}

const del = (e) => {
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
    
    getParentRow(e).remove() // Delete entire row

    return false; // Don't postback
}
