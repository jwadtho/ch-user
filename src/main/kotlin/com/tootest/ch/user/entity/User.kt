package com.tootest.ch.user.entity

import javax.persistence.Entity
import javax.persistence.GeneratedValue
import javax.persistence.GenerationType
import javax.persistence.Id

@Entity
class User(
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    var id: Int? = null,
    var name: String,
    var email: String
) {
    override fun toString(): String {
        return "User(id=$id, name='$name', email='$email')"
    }
}