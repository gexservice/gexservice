//auto gen func by autogen
package basedb

/**
 * @apiDefine AnnounceUpdate
 * @apiParam (Announce) {String} Announce.title only required when add, the announce title
 * @apiParam (Announce) {Int} [Announce.marked] the announce marked
 * @apiParam (Announce) {Object} [Announce.info] the announce external info
 * @apiParam (Announce) {Object} [Announce.content] the announce content
 * @apiParam (Announce) {AnnounceStatus} [Announce.status] the announce status, all suported is <a href="#metadata-Announce">AnnounceStatusAll</a>
 */
/**
 * @apiDefine AnnounceObject
 * @apiSuccess (Announce) {Int64} Announce.tid the announce id
 * @apiSuccess (Announce) {AnnounceType} Announce.type the announce type, all suported is <a href="#metadata-Announce">AnnounceTypeAll</a>
 * @apiSuccess (Announce) {Int} Announce.marked the announce marked
 * @apiSuccess (Announce) {String} Announce.title the announce title
 * @apiSuccess (Announce) {Object} Announce.info the announce external info
 * @apiSuccess (Announce) {Object} Announce.content the announce content
 * @apiSuccess (Announce) {Time} Announce.update_time the announce update time
 * @apiSuccess (Announce) {Time} Announce.create_time the announce create time
 * @apiSuccess (Announce) {AnnounceStatus} Announce.status the announce status, all suported is <a href="#metadata-Announce">AnnounceStatusAll</a>
 */

/**
 * @apiDefine ConfigUpdate
 */
/**
 * @apiDefine ConfigObject
 * @apiSuccess (Config) {String} Config.key
 * @apiSuccess (Config) {String} Config.value
 * @apiSuccess (Config) {Time} Config.update_time
 */

/**
 * @apiDefine ObjectUpdate
 */
/**
 * @apiDefine ObjectObject
 * @apiSuccess (Object) {String} Object.key the object key
 * @apiSuccess (Object) {Object} Object.value the object value
 * @apiSuccess (Object) {Time} Object.update_time
 * @apiSuccess (Object) {Time} Object.create_time the create time
 * @apiSuccess (Object) {ObjectStatus} Object.status the status, all suported is <a href="#metadata-Object">ObjectStatusAll</a>
 */

/**
 * @apiDefine VersionObjectUpdate
 * @apiParam (VersionObject) {String} VersionObject.pub only required when add, the publish scoe of version object,split multi by comma,* to all,x.x.x.x for ip
 * @apiParam (VersionObject) {Object} VersionObject.value only required when add, the version of key
 * @apiParam (VersionObject) {VersionObjectStatus} VersionObject.status only required when add, the status, all suported is <a href="#metadata-VersionObject">VersionObjectStatusAll</a>
 */
/**
 * @apiDefine VersionObjectObject
 * @apiSuccess (VersionObject) {Int64} VersionObject.tid the primary key
 * @apiSuccess (VersionObject) {String} VersionObject.key the name of key
 * @apiSuccess (VersionObject) {String} VersionObject.pub the publish scoe of version object,split multi by comma,* to all,x.x.x.x for ip
 * @apiSuccess (VersionObject) {Object} VersionObject.value the version of key
 * @apiSuccess (VersionObject) {Time} VersionObject.update_time the update time
 * @apiSuccess (VersionObject) {Time} VersionObject.create_time the create time
 * @apiSuccess (VersionObject) {VersionObjectStatus} VersionObject.status the status, all suported is <a href="#metadata-VersionObject">VersionObjectStatusAll</a>
 */
